package database

import (
    "net/http"
    "expvar"
    "strconv"
    "strings"
    "encoding/json"

    "github.com/jsirianni/systemstat/internal/log"
    "github.com/jsirianni/systemstat/internal/email"

    "github.com/gorilla/mux"
)

// counter metrics exposed at /debug/vars
var counts = expvar.NewMap("counters")

type Server struct {
    Port int
    DB   Database
}

func init() {
    counts.Add("total_requests", 0)
}

func (s Server) Run() error {
    port := strconv.Itoa(s.Port)

    log.Info("starting database api on port:", port)

    router := mux.NewRouter()
    router.HandleFunc("/status", s.status).Methods("GET")
    router.HandleFunc("/v1/account/{email}", s.getAccount).Methods("GET")
    router.HandleFunc("/v1/account/{email}", s.createAccount).Methods("POST")
    // expvar runtime  metrics
    router.Handle("/debug/vars", http.DefaultServeMux)
    return http.ListenAndServe(":" + port, router)
}

func (s Server) status(resp http.ResponseWriter, req *http.Request) {
	if err := s.DB.TestConnection(); err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
        return
	}
	resp.WriteHeader(http.StatusOK)
}

func (s Server) createAccount(resp http.ResponseWriter, req *http.Request) {
    counts.Add("total_requests", 1)

    emailAddr :=  mux.Vars(req)["email"]
    if emailAddr == "" {
        resp.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := email.Validate(emailAddr); err != nil {
        http.Error(resp, "email address is not valid", http.StatusUnprocessableEntity)
        return
    }

    if _, err := s.DB.AccountByEmail(emailAddr); err == nil {
        resp.WriteHeader(http.StatusConflict)
        return
    }

    account, err := s.DB.AccountCreate(emailAddr)
    if err != nil {
        log.Error(err)
        resp.WriteHeader(http.StatusInternalServerError)
        return
    }

    resp.Header().Set("Content-Type", "application/json")
    resp.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(resp).Encode(account); err != nil {
        log.Error(err)
        return
    }
}

func (s Server) getAccount(resp http.ResponseWriter, req *http.Request) {
    counts.Add("total_requests", 1)

    emailAddr :=  mux.Vars(req)["email"]
    if emailAddr == "" {
        resp.WriteHeader(http.StatusBadRequest)
        return
    }

    account, err := s.DB.AccountByEmail(emailAddr)
    if err != nil {
        log.Error(err)
        resp.WriteHeader(http.StatusNotFound)
        return
    }

    resp.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(resp).Encode(account); err != nil {
        log.Error(err)
        return
    }
}

func isApplicationJSON(req *http.Request) bool {
    if strings.ToLower(req.Header.Get("Content-Type")) != "application/json" {
        return false
    }
    return true
}
