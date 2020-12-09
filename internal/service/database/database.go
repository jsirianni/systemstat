package database

import (
	"github.com/jsirianni/systemstat/internal/service/database/postgres"
)

type Database interface {
	Validate() error
	TestConnection() error
	Insert() error
	Select() (string, error)
}

func NewPostgres() (Database, error) {
	d, err := postgres.New()
	if err != nil {
		return nil, err
	}
	return d, d.Validate()
}
