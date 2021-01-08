package database

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jsirianni/systemstat/internal/email"
	"github.com/jsirianni/systemstat/internal/log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Port struct {
		GRPC int
	}
	DB Database

	UnimplementedApiServer
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

func (s Server) HealthCheck(c context.Context, req *HealthRequest) (*Health, error) {
	h, err := s.health()
	return &h, err
}

func (s Server) GetAccount(c context.Context, req *GetAccountRequest) (*Account, error) {
	a, err := s.getAccount(req.AccountId)
	return &a, err
}

func (s Server) CreateToken(ctx context.Context, req *CreateTokenRequest) (*Token, error) {
	t, err := s.createToken()
	return &t, err
}

func (s Server) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	account, err := s.createAccount(req.Email, req.Token)
	return &account, err
}

func (s Server) health() (h Health, err error) {
	err = s.DB.TestConnection()
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s Server) getAccount(id string) (Account, error) {
	a, err := s.DB.AccountByID(id)
	if err != nil {
		log.Debug(err)
		return a, err
	}
	log.Trace("getAccount: account retrieved:", id, a.AdminEmail)
	return a, nil
}

func (s Server) createToken() (Token, error) {
	t, err := s.DB.CreateToken()
	if err != nil {
		log.Error(err)
		return t, err
	}
	log.Trace("createToken: token created:", t.Token)
	return t, nil
}

func (s Server) createAccount(emailAddr, token string) (account Account, err error) {
	log.Trace(fmt.Sprintf("createAccount: requested creation: email=%s token=%s", emailAddr, token))

	const invalidtoken = "signup token is invalid"

	if err := email.Validate(emailAddr); err != nil {
		log.Debug(err)
		return account, err
	}

	// check if account exists first, err will not be nil if the account
	// does not exist
	if a, err := s.DB.AccountByEmail(emailAddr); err == nil {
		err := errors.New(fmt.Sprintf("account with email address %s already exists with account_id %s", emailAddr, a.AccountId))
		log.Debug(err)
		return account, err
	}

	// claim the token before creating the account
	if _, err := s.DB.ClaimToken(emailAddr, token); err != nil {
		log.Error(err)
		return account, err
	}

	account, err = s.DB.AccountCreate(emailAddr, token)
	if err != nil {
		log.Error(err)
		return account, err
	}
	log.Trace("createAccount: account created: ", account.AccountId, account.AdminEmail)
	return account, nil
}
