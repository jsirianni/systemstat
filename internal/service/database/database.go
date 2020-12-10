package database

import (
	"github.com/jsirianni/systemstat/internal/service/database/postgres"
	"github.com/jsirianni/systemstat/internal/types/account"
)

type Database interface {
	Validate() error
	TestConnection() error

	AccountCreate(email string) (account.Account, error)
	AccountByEmail(email string) (account.Account, error)
	AccountConfigureAlert(alertType string, config []byte) (account.Account, error)
}

func NewPostgres() (Database, error) {
	d, err := postgres.New()
	if err != nil {
		return nil, err
	}
	return d, d.Validate()
}
