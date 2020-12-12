package database

import (
    "os"
    "fmt"
    "net"
    "strconv"
    "context"

    "github.com/jsirianni/systemstat/internal/log"
    "github.com/jsirianni/systemstat/api"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
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
	a, err := s.getAccount(req.AccountId)
    return &a, err
}

func (s Server) CreateToken(ctx context.Context, req *api.CreateTokenRequest) (*api.Token, error) {
    t, err := s.createToken()
    return &t, err
}

func (s Server) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.Account, error) {
	account, err := s.createAccount(req.Email, req.Token)
    return &account, err
}
