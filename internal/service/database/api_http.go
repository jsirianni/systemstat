package database

import (
	"encoding/json"
	"expvar"
	"net/http"
	"strconv"

	"github.com/jsirianni/systemstat/internal/email"
	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/proto"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// counter metrics exposed at /debug/vars
var counts = expvar.NewMap("counters")

type Server struct {
	Port struct {
		HTTP int
		GRPC int
	}
	DB   Database

	proto.UnimplementedApiServer
}

func init() {
	counts.Add("total_requests", 0)
}

func (s Server) RunHTTP() error {
	port := strconv.Itoa(s.Port.HTTP)

	log.Info("starting http api on port:", port)

	router := mux.NewRouter()
	router.HandleFunc("/status", s.status).Methods("GET")
	router.HandleFunc("/v1/account/token/create", s.createTokenHandler).Methods("POST")
	router.HandleFunc("/v1/account/{account_id}", s.getAccountHandler).Methods("GET")
	router.HandleFunc("/v1/account/{token}/{email}", s.createAccountHandler).Methods("POST")
	// expvar runtime  metrics
	router.Handle("/debug/vars", http.DefaultServeMux)
	return http.ListenAndServe(":"+port, router)
}

func (s Server) status(resp http.ResponseWriter, req *http.Request) {
	if err := s.DB.TestConnection(); err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (s Server) createTokenHandler(resp http.ResponseWriter, req *http.Request) {
	t, err := s.DB.CreateToken()
	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Trace("token created:", t.Token.String())

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(resp).Encode(t); err != nil {
		log.Error(err)
		return
	}
}

func (s Server) createAccountHandler(resp http.ResponseWriter, req *http.Request) {
	counts.Add("total_requests", 1)

	emailAddr := mux.Vars(req)["email"]
	token := mux.Vars(req)["token"]
	if emailAddr == "" || token == "" {
		log.Debug(errors.New("createAccount: client request missing 'email' or 'token' parameter"))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	const invalidtoken = "signup token is invalid"

	if err := email.Validate(emailAddr); err != nil {
		log.Debug(err)
		http.Error(resp, "email addess is not valid", http.StatusUnprocessableEntity)
		return
	}

	if _, err := uuid.Parse(token); err != nil {
		log.Debug(err)
		http.Error(resp, invalidtoken, http.StatusUnprocessableEntity)
		return
	}

	// check if account exists first, err will not be nil if the account
	// does not exist
	if _, err := s.DB.AccountByEmail(emailAddr); err == nil {
		log.Debug(errors.New("account already exists"))
		resp.WriteHeader(http.StatusConflict)
		return
	}

	// claim the token before creating the account
	if _, err := s.DB.ClaimToken(emailAddr, token); err != nil {
		log.Debug(err)
		http.Error(resp, invalidtoken, http.StatusUnprocessableEntity)
		return
	}

	account, err := s.DB.AccountCreate(emailAddr, token)
	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Trace("account created:", emailAddr)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(resp).Encode(account); err != nil {
		log.Error(err)
		return
	}
}

func (s Server) getAccountHandler(resp http.ResponseWriter, req *http.Request) {
	counts.Add("total_requests", 1)

	id := mux.Vars(req)["account_id"]
	if id == "" {
		log.Debug(errors.New("getAccount: client request missing 'account_id' variable"))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	genericErrContext := errors.New("getAccount: account_id " + id)

	account, err := s.DB.AccountByID(id)
	if err != nil {
		log.Debug(errors.Wrap(err, genericErrContext.Error()))
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(account); err != nil {
		log.Error(errors.Wrap(err, genericErrContext.Error()))
		return
	}
	log.Trace("getAccount: account retrieved:", id, account.AdminEmail)
}
