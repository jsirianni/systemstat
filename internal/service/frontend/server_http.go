package frontend

import (
	"fmt"
	"net/http"
	"strconv"
	"context"
	"encoding/json"

	"github.com/jsirianni/systemstat/internal/log"
	"github.com/jsirianni/systemstat/api"

	"github.com/gorilla/mux"
)

func (s Server) httpServer() error {
	router := mux.NewRouter()

	// health endpoint returns status ok
	router.HandleFunc("/health", s.status).Methods("GET")

	// get account endpoint returns an api.Account
	// takes an api key as header x-api-key
	// takes an account id in query string as 'account'
	router.HandleFunc("/v1/account",s.httpGetAccount).Queries("account", "{[0-9]*?}").Methods("GET").Headers(headerAPIKey, "")

	// create token endpoint creates and returns a new and unclaimed signup api.Token
	// takes an api key (admin only) as header x-api-key
	router.HandleFunc("/v1/token",s.httpCreateToken).Methods("POST").Headers(headerAPIKey, "")

	// create account endpoint returns a new api.Account
	// takes a token and email in the payload
	router.HandleFunc("/v1/account",s.httpCreateAccount).Methods("POST")

	port := strconv.Itoa(s.Port.HTTP)
	log.Info("starting frontend api on port:", port)
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

func (s Server) httpGetAccount(resp http.ResponseWriter, req *http.Request) {
	r := api.GetAccountRequest{
		ApiKey: req.Header.Get(headerAPIKey),
		AccountId: req.URL.Query().Get("account"),
	}

	account, err := s.GetAccount(context.Background(), &r)
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

func (s Server) httpCreateToken(resp http.ResponseWriter, req *http.Request) {
	r := api.CreateTokenRequest{
		ApiKey: req.Header.Get(headerAPIKey),
	}

	token, err := s.CreateToken(context.Background(), &r)
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

func (s Server) httpCreateAccount(resp http.ResponseWriter, req *http.Request) {
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
