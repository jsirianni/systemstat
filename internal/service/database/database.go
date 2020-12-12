package database

import (
	"github.com/jsirianni/systemstat/internal/service/database/postgres"
	"github.com/jsirianni/systemstat/internal/proto"
	"github.com/jsirianni/systemstat/internal/types/account"
)

type Database interface {
	// Validate the database configuration
	Validate() error

	// Test connection to the database
	TestConnection() error

	// Create an account
	AccountCreate(email, token string) (proto.Account, error)

	// Retrieve an account by account_id
	AccountByID(id string) (proto.Account, error)

	// Retrieve an account by email
	AccountByEmail(email string) (proto.Account, error)

	// Configure an accounts alert type
	AccountConfigureAlert(alertType string, config []byte) (proto.Account, error)

	// claim a sign up token
	ClaimToken(email, token string) (account.Token, error)

	// get an existing token
	GetToken(token string) (account.Token, error)

	// create a signup token
	CreateToken() (account.Token, error)
}

func NewPostgres() (Database, error) {
	d, err := postgres.New()
	if err != nil {
		return nil, err
	}
	return d, d.Validate()
}
