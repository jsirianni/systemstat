package database

import (
	"github.com/jsirianni/systemstat/internal/service/database/api"
	"github.com/jsirianni/systemstat/internal/service/database/postgres"
)

type Database interface {
	// Validate the database configuration
	Validate() error

	// Test connection to the database
	TestConnection() error

	// Create an account
	AccountCreate(email, token string) (api.Account, error)

	// Retrieve an account by account_id
	AccountByID(id string) (api.Account, error)

	// Retrieve an account by email
	AccountByEmail(email string) (api.Account, error)

	// Configure an accounts alert type
	AccountConfigureAlert(alertType string, config []byte) (api.Account, error)

	// claim a sign up token
	ClaimToken(email, token string) (api.Token, error)

	// get an existing token
	GetToken(token string) (api.Token, error)

	// create a signup token
	CreateToken() (api.Token, error)
}

func NewPostgres() (Database, error) {
	d, err := postgres.New()
	if err != nil {
		return nil, err
	}
	return d, d.Validate()
}
