package postgres

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"
)

type Postgres struct {
	host   string
	port   int
	user   string
	pass   string
	dbname string

	db *sql.DB
}

func New() (Postgres, error) {
	p := Postgres{}
	p.configure()
	if err := p.Validate(); err != nil {
		return p, err
	}

	var conn = p.connectionString()
	var err error

	p.db, err = sql.Open("postgres", conn)
	if err != nil {
		return p, err
	}
	return p, p.TestConnection()
}

func (p Postgres) TestConnection() error {
	return p.db.Ping()
}

func (p Postgres) connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.host, p.port, p.user, p.pass, p.dbname)
}
