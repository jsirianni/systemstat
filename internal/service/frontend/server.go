package frontend

import (
    "os"
    "fmt"
    "net"
    "strconv"
    "context"

    "github.com/jsirianni/systemstat/internal/service/database"
    "github.com/jsirianni/systemstat/internal/log"

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

	database.RegisterApiServer(grpcServer, s)

	log.Info("starting grpc api on port:", strconv.Itoa(s.Port.GRPC))
	return grpcServer.Serve(lis)
}

func (s Server) HealthCheck(c context.Context, req *database.HealthRequest) (*database.Health, error) {
    return nil, nil
}

func (s Server) GetAccount(context.Context, *database.GetAccountRequest) (*database.Account, error) {
    return nil, nil
}

func (s Server) CreateToken(context.Context, *database.CreateTokenRequest) (*database.Token, error) {
    return nil, nil
}

func (s Server) CreateAccount(context.Context, *database.CreateAccountRequest) (*database.Account, error) {
    return nil, nil
}

func (s Server) authenticated(apiKey string, admin bool) bool {
    // TODO: use admin bool

    if apiKey == os.Getenv(envAdminToken) {
        return true
    }
    return false
}
