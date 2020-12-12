package database

import (
    "os"
    "fmt"
    "net"
    "strconv"
    "context"

    "github.com/jsirianni/systemstat/internal/email"
    "github.com/jsirianni/systemstat/internal/log"
    "github.com/jsirianni/systemstat/api"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/pkg/errors"
    "github.com/google/uuid"
)

func (s Server) RunGRPC() error {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port.GRPC))
    if err != nil {
      log.Fatal(err, 1)
    }

    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)

    // allow grpcurl https://github.com/fullstorydev/grpcurl
    if os.Getenv("GO_ENV") == "development" {
        log.Trace("GO_ENV=development detected, enabling GRPC reflection")
        reflection.Register(grpcServer)
    }

    api.RegisterApiServer(grpcServer, s)

    log.Info("starting grpc api on port:", strconv.Itoa(s.Port.GRPC))
    return grpcServer.Serve(lis)
}

func (s Server) GetAccount(c context.Context, req *api.GetAccountRequest) (*api.Account, error) {
	a, err := s.DB.AccountByID(req.AccountId)
    return &a, err
}

func (s Server) CreateToken(ctx context.Context, req *api.CreateTokenRequest) (*api.Token, error) {
    t, err := s.DB.CreateToken()
    return &t, err
}

func (s Server) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.Account, error) {
    emailAddr := req.Email
	token := req.Token
	if emailAddr == "" || token == "" {
		return &api.Account{}, errors.New("createAccount: client request missing 'email' or 'token' parameter")
	}

	const invalidtoken = "signup token is invalid"

	if err := email.Validate(emailAddr); err != nil {
		return &api.Account{}, err
	}

	if _, err := uuid.Parse(token); err != nil {
		return &api.Account{}, err
	}

	// check if account exists first, err will not be nil if the account
	// does not exist
	if _, err := s.DB.AccountByEmail(emailAddr); err == nil {
		return &api.Account{}, err
	}

	// claim the token before creating the account
	if _, err := s.DB.ClaimToken(emailAddr, token); err != nil {
		return &api.Account{}, err
	}

	account, err := s.DB.AccountCreate(emailAddr, token)
    return &account, err
}
