package postgres

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	envHost   = "DATABASE_HOST"
	envPort   = "DATABASE_PORT"
	envUser   = "DATABASE_USERNAME"
	envPass   = "DATABASE_PASSWORD"
	envDBName = "DATABASE_NAME"
)

func (p *Postgres) configure() error {
	p.host = os.Getenv(envHost)
	p.user = os.Getenv(envUser)
	p.pass = os.Getenv(envPass)
	p.dbname = os.Getenv(envDBName)
	return p.configurePort()
}

func (p *Postgres) configurePort() error {
	p.port = 5432
	if port := os.Getenv(envPort); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return errors.Wrap(err, envPort+" is not an int, got "+port)
		}
		p.port = portInt
	}
	return nil
}
