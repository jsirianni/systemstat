package database

import (
    "os"
    "fmt"
    "net"
    "strconv"
    "context"


    "github.com/jsirianni/systemstat/internal/log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

func (s Server) GetAccount(c context.Context, req *GetAccountRequest) (*GetAccountReply, error) {
	a, err := s.DB.AccountByID(req.AccountId)
	if err != nil {
		return &GetAccountReply{}, nil
	}

	acct := GetAccountReply{}
	acct.AccountId = a.AccountID.String()
	acct.RootApiKey = a.RootAPIKey.String()
	acct.AlertType = a.AlertType
	acct.AdminEmail = a.AdminEmail
	acct.AlertConfig, err = a.AlertConfig.JSON()
	return &acct, err
}

func (s Server) mustEmbedUnimplementedApiServer() {
	return
}

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

    RegisterApiServer(grpcServer, s)

    log.Info("starting grpc api on port:", strconv.Itoa(s.Port.GRPC))
    return grpcServer.Serve(lis)
}
