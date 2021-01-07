package frontend

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jsirianni/systemstat/internal/log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	Port     int
	Database struct {
		Endpoint string
	}
}

func (s Server) Run() error {
	port := strconv.Itoa(s.Port)

	log.Info("starting frontend api on port:", port)

	router := mux.NewRouter()
	router.HandleFunc("/status", s.status).Methods("GET")
	// expvar runtime  metrics
	router.Handle("/debug/vars", http.DefaultServeMux)
	return http.ListenAndServe(":"+port, router)
}

func (s Server) status(resp http.ResponseWriter, req *http.Request) {
	backendStatus, err := http.Get(s.Database.Endpoint + "/status")
	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	status := strconv.Itoa(backendStatus.StatusCode)
	if backendStatus.StatusCode != http.StatusOK {
		log.Error(errors.New(fmt.Sprintf("backend returned status %s", status)))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
}
