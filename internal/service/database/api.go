package database

import (
    "fmt"
    "net/http"

    "github.com/jsirianni/systemstat/api"
    "github.com/jsirianni/systemstat/internal/log"
    "github.com/jsirianni/systemstat/internal/email"

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

func (s Server) status() (h api.Health, err error) {
    err = s.DB.TestConnection()
    if err != nil {
        log.Error(err)
        h.HttpStatus = errorToHTTPStatus(err)
        return
    }
    h.HttpStatus = http.StatusOK
    return
}

func (s Server) getAccount(id string) (api.Account, error) {
    a, err := s.DB.AccountByID(id)
    if err != nil {
        log.Debug(err)
        a.HttpStatus = errorToHTTPStatus(err)
        return a, err
    }
    log.Trace("getAccount: account retrieved:", id, a.AdminEmail)
    a.HttpStatus = http.StatusOK
    return a, nil
}

func (s Server) createToken() (api.Token, error) {
    t, err := s.DB.CreateToken()
    if err != nil {
        log.Error(err)
        t.HttpStatus = errorToHTTPStatus(err)
        return t, err
    }
    log.Trace("createToken: token created:", t.Token)
    t.HttpStatus = http.StatusCreated
    return t, nil
}

func (s Server) createAccount(emailAddr, token string) (account api.Account, err error) {
    log.Trace(fmt.Sprintf("createAccount: requested creation: email=%s token=%s", emailAddr, token))

    const invalidtoken = "signup token is invalid"

    if err := email.Validate(emailAddr); err != nil {
        log.Debug(err)
        account.HttpStatus = http.StatusBadRequest
        return account, err
    }

    // check if account exists first, err will not be nil if the account
    // does not exist
    if a, err := s.DB.AccountByEmail(emailAddr); err == nil {
        log.Debug(errors.New(fmt.Sprintf("account with email address %s already exists with account_id %s", emailAddr, a.AccountId)))
        account.HttpStatus = http.StatusConflict
        return account, err
    }

    // claim the token before creating the account
    if _, err := s.DB.ClaimToken(emailAddr, token); err != nil {
        log.Error(err)
        account.HttpStatus = errorToHTTPStatus(err)
        return account, err
    }

    account, err = s.DB.AccountCreate(emailAddr, token)
    if err != nil {
        log.Error(err)
        account.HttpStatus = errorToHTTPStatus(err)
        return account, err
    }
    log.Trace("createAccount: account created: ", account.AccountId, account.AdminEmail)
    account.HttpStatus = http.StatusCreated
    return account, nil
}
