package database

import (
    "fmt"

    "github.com/jsirianni/systemstat/api"
    "github.com/jsirianni/systemstat/internal/log"
    "github.com/jsirianni/systemstat/internal/email"

    "github.com/google/uuid"
    "github.com/pkg/errors"
)

type Server struct {
	Port struct {
		HTTP int
		GRPC int
	}
	DB   Database

	api.UnimplementedApiServer
}

func (s Server) getAccount(id string) (api.Account, error) {
    a, err := s.DB.AccountByID(id)
    if err != nil {
        log.Debug(err)
        return a, err
    }
    log.Trace("getAccount: account retrieved:", id, a.AdminEmail)
    return a, nil
}

func (s Server) createToken() (api.Token, error) {
    t, err := s.DB.CreateToken()
    if err != nil {
        log.Error(err)
        return t, err
    }
    log.Trace("createToken: token created:", t.Token)
    return t, nil
}

func (s Server) createAccount(emailAddr, token string) (api.Account, error) {
    log.Trace(fmt.Sprintf("createAccount: requested creation: email=%s token=%s", emailAddr, token))

    const invalidtoken = "signup token is invalid"

    if err := email.Validate(emailAddr); err != nil {
        log.Debug(err)
        return api.Account{}, err
    }

    if _, err := uuid.Parse(token); err != nil {
        log.Debug(err)
        return api.Account{}, err
    }

    // check if account exists first, err will not be nil if the account
    // does not exist
    if _, err := s.DB.AccountByEmail(emailAddr); err == nil {
        err := errors.New(fmt.Sprintf("account with email %s already exists", emailAddr))
        log.Debug(err)
        return api.Account{}, err
    }

    // claim the token before creating the account
    if _, err := s.DB.ClaimToken(emailAddr, token); err != nil {
        log.Error(err)
        return api.Account{}, err
    }

    account, err := s.DB.AccountCreate(emailAddr, token)
    if err != nil {
        log.Error(err)
        return api.Account{}, err
    }
    log.Trace("createAccount: account created: ", account.AccountId, account.AdminEmail)
    return account, nil
}
