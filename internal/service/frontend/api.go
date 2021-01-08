package frontend

import (
	"os"
	"fmt"
	"net/http"
	"strconv"
	"context"
	"encoding/json"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/internal/service/database/api"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const envAdminToken = "SYSTEMSTAT_ADMIN_TOKEN"
const headerAPIKey  = "X-Api-Key"

type Server struct {
	Port     int

	Database struct {
		// GRCP endpoint
		Endpoint string

		rpcConn *grpc.ClientConn
		client api.ApiClient
	}
}

func (s Server) Run() error {
	if err := s.initServer(); err != nil {
		return errors.Wrap(err, "failed to initialize frontend api")
	}

	port := strconv.Itoa(s.Port)

	var err error
	s.Database.rpcConn, err = grpc.Dial(s.Database.Endpoint, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to initialize backend database grpc client")
	}
	defer s.Database.rpcConn.Close()

	s.Database.client = api.NewApiClient(s.Database.rpcConn)

	log.Info("starting frontend api on port:", port)

	router := mux.NewRouter()

	// health endpoint returns status ok
	router.HandleFunc("/health", s.status).Methods("GET")

	// get account endpoint returns an api.Account
	// takes an api key as header x-api-key
	// takes an account id in query string as 'account'
	router.HandleFunc("/v1/account",s.getAccount).Queries("account", "{[0-9]*?}").Methods("GET").Headers(headerAPIKey, "")

	// create token endpoint creates and returns a new and unclaimed signup api.Token
	// takes an api key (admin only) as header x-api-key
	router.HandleFunc("/v1/token",s.createToken).Methods("POST").Headers(headerAPIKey, "")

	// create account endpoint returns a new api.Account
	// takes a token and email in the payload
	router.HandleFunc("/v1/account",s.createAccount).Methods("POST")

	return http.ListenAndServe(":"+port, router)
}

func (s Server) status(resp http.ResponseWriter, req *http.Request) {
	_, err := s.Database.client.HealthCheck(context.Background(), &api.HealthRequest{})
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (s Server) getAccount(resp http.ResponseWriter, req *http.Request) {
 	if ! s.authenticated(req, false) {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	r := api.GetAccountRequest{
		AccountId: req.URL.Query().Get("account"),
	}

	account, err := s.Database.client.GetAccount(context.Background(), &r)
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, string(b))
}

func (s Server) createToken(resp http.ResponseWriter, req *http.Request) {
	if ! s.authenticated(req, true) {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := s.Database.client.CreateToken(context.Background(), &api.CreateTokenRequest{})
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	fmt.Fprint(resp, string(b))
}

func (s Server) createAccount(resp http.ResponseWriter, req *http.Request) {
	r := api.CreateAccountRequest{}
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	account, err := s.Database.client.CreateAccount(context.Background(), &r)
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		log.Debug(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	fmt.Fprint(resp, string(b))
}

func (s Server) authenticated(req *http.Request, admin bool) bool {
    // TODO: use admin bool

	apiKey := req.Header.Get(headerAPIKey)
	if apiKey == "" {
		log.Debug(fmt.Errorf("%q header is missing", headerAPIKey))
		return false
	}

    if apiKey == os.Getenv(envAdminToken) {
        return true
    }
    return false
}


func (s *Server) initServer() error {
	if s.Port == 0 {
		return errors.New("port not set")
	}

	if s.Database.Endpoint == "" {
		return errors.New("database endpoint not set")
	}

	return nil
}
