package frontend

import (
    "os"
    "fmt"
    "net"
    "strconv"
    "context"

    "github.com/jsirianni/systemstat/api"
    "github.com/jsirianni/systemstat/internal/log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/pkg/errors"
)

const envAdminToken = "SYSTEMSTAT_ADMIN_TOKEN"
const headerAPIKey  = "X-Api-Key"

type Server struct {
	Port struct {
		HTTP int
		GRPC int
	}

	Database struct {
		// GRCP endpoint
		Endpoint string

		rpcConn *grpc.ClientConn
		client api.DatabaseClient
	}

	api.UnimplementedFrontendServer
}

func (s *Server) Server() error {
    // connect to backend database grpc server
    var err error
    s.Database.rpcConn, err = grpc.Dial(s.Database.Endpoint, grpc.WithInsecure())
    if err != nil {
        return errors.Wrap(err, "failed to initialize backend database grpc client")
    }
    defer s.Database.rpcConn.Close()
    s.Database.client = api.NewDatabaseClient(s.Database.rpcConn)

    // start frontend grpc server
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port.GRPC))
    if err != nil {
        return err
    }

    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)

    // allow grpcurl https://github.com/fullstorydev/grpcurl
    if os.Getenv("GO_ENV") == "development" {
        log.Trace("GO_ENV=development detected, enabling GRPC reflection")
        reflection.Register(grpcServer)
    }

    api.RegisterFrontendServer(grpcServer, s)

    go s.httpServer()

    log.Info("starting grpc api on port:", strconv.Itoa(s.Port.GRPC))
	return grpcServer.Serve(lis)
}

func (s Server) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.Account, error) {
    if ! s.authenticated(req.ApiKey, false) {
        return nil, fmt.Errorf("api key is required")
    }

	return s.Database.client.GetAccount(ctx, req)
}

func (s Server) CreateToken(ctx context.Context, req *api.CreateTokenRequest) (*api.Token, error) {
    if ! s.authenticated(req.ApiKey, false) {
        return nil, fmt.Errorf("api key is required")
    }

    return s.Database.client.CreateToken(ctx, req)
}

func (s Server) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.Account, error) {
    return s.Database.client.CreateAccount(ctx, req)
}

func (s Server) authenticated(apiKey string, admin bool) bool {
    // TODO: use admin bool
    // NOTE: right now, we treat all api key's as if they
    // are admin.

    if apiKey == os.Getenv(envAdminToken) {
        return true
    }
    return false
}
