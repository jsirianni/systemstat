package database

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jsirianni/systemstat/internal/log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s Server) RunHTTP() error {
	port := strconv.Itoa(s.Port.HTTP)

	log.Info("starting http api on port:", port)

	router := mux.NewRouter()
	router.HandleFunc("/status", s.statusHandler).Methods("GET")
	router.HandleFunc("/v1/account/token/create", s.createTokenHandler).Methods("POST")
	router.HandleFunc("/v1/account/{account_id}", s.getAccountHandler).Methods("GET")
	router.HandleFunc("/v1/account/{token}/{email}", s.createAccountHandler).Methods("POST")
	return http.ListenAndServe(":"+port, router)
}

func (s Server) statusHandler(resp http.ResponseWriter, req *http.Request) {
	h, err := s.status()
	if err != nil {
		resp.WriteHeader(int(h.HttpStatus))
		return
	}
	resp.WriteHeader(int(h.HttpStatus))
}

func (s Server) createTokenHandler(resp http.ResponseWriter, req *http.Request) {
	t, err := s.createToken()
	if err != nil {
		resp.WriteHeader(int(t.HttpStatus))
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(int(t.HttpStatus))
	if err := json.NewEncoder(resp).Encode(t); err != nil {
		log.Error(err)
		return
	}
}

func (s Server) createAccountHandler(resp http.ResponseWriter, req *http.Request) {
	emailAddr := mux.Vars(req)["email"]
	token := mux.Vars(req)["token"]
	if emailAddr == "" || token == "" {
		log.Debug(errors.New("createAccount: client request missing 'email' or 'token' parameter"))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	account, err := s.createAccount(emailAddr, token)
	if err != nil {
		resp.WriteHeader(int(account.HttpStatus))
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(int(account.HttpStatus))
	if err := json.NewEncoder(resp).Encode(account); err != nil {
		log.Error(err)
		return
	}
}

func (s Server) getAccountHandler(resp http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["account_id"]
	if id == "" {
		log.Debug(errors.New("getAccount: client request missing 'account_id' variable"))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := s.getAccount(id)
	if err != nil {
		resp.WriteHeader(int(account.HttpStatus))
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(int(account.HttpStatus))
	if err := json.NewEncoder(resp).Encode(account); err != nil {
		log.Error(err)
		return
	}
}
