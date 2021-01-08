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

	// account endpoint returns an api.Account
	// takes an api key as header x-api-key
	// takes an account id in query string as 'account'
	router.HandleFunc("/v1/account",s.getAccount).Queries("account", "{[0-9]*?}").Methods("GET").Headers(headerAPIKey, "")

	// token endpoint creates and returns a new and unclaimed signup api.Token
	// takes an api key (admin only) as header x-api-key
	router.HandleFunc("/v1/token",s.createToken).Methods("POST").Headers(headerAPIKey, "")

	return http.ListenAndServe(":"+port, router)
}

func (s Server) status(resp http.ResponseWriter, req *http.Request) {
	_, err := s.Database.client.HealthCheck(context.Background(), &api.HealthRequest{})
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (s Server) getAccount(resp http.ResponseWriter, req *http.Request) {
	a, err := s.authenticated(req, false)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	} else if a == false {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	r := api.GetAccountRequest{
		AccountId: req.URL.Query().Get("account"),
	}

	account, err := s.Database.client.GetAccount(context.Background(), &r)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(account)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(resp, string(b))
	resp.WriteHeader(http.StatusOK)
}

func (s Server) createToken(resp http.ResponseWriter, req *http.Request) {
	a, err := s.authenticated(req, true)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	} else if a == false {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := s.Database.client.CreateToken(context.Background(), &api.CreateTokenRequest{})
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(resp, string(b))
	resp.WriteHeader(http.StatusCreated)
}

func (s Server) authenticated(req *http.Request, admin bool) (bool, error) {
    // TODO: use admin bool

	apiKey := req.Header.Get(headerAPIKey)
	if apiKey == "" {
		return false, fmt.Errorf("%q header is missing", headerAPIKey)
	}

    if apiKey == os.Getenv(envAdminToken) {
        return true, nil
    }
    return false, nil
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
