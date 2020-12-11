package postgres

import (
	"fmt"

	"github.com/jsirianni/systemstat/internal/types/account"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/pkg/errors"
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

	var conn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.host, p.port, p.user, p.pass, p.dbname)
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

func (p Postgres) AccountCreate(email string) (account.Account, error) {
	if email == "" {
		return account.Account{}, errors.New("email is a required parameter when creating an account")
	}

	q := fmt.Sprintf("INSERT INTO account (admin_email) VALUES ('%s')", email)
	if _, err := p.db.Exec(q); err != nil {
		return account.Account{}, err
	}

	return p.AccountByEmail(email)
}

func (p Postgres) AccountByID(id string) (account.Account, error) {
	a := account.Account{}

	if id == "" {
		return a, errors.New("account_id is a required parameter when reading an account")
	}

	q := fmt.Sprintf("SELECT * FROM account WHERE account_id = '%s'", id)
	err := p.db.QueryRow(q).Scan(&a.AccountID, &a.RootAPIKey, &a.AlertType, &a.AlertConfig, &a.AdminEmail)
	return a, err
}

func (p Postgres) AccountByEmail(email string) (account.Account, error) {
	a := account.Account{}

	if email == "" {
		return a, errors.New("email is a required parameter when reading an account")
	}

	q := fmt.Sprintf("SELECT * FROM account WHERE admin_email = '%s'", email)
	err := p.db.QueryRow(q).Scan(&a.AccountID, &a.RootAPIKey, &a.AlertType, &a.AlertConfig, &a.AdminEmail)
	return a, err
}

func (p Postgres) AccountConfigureAlert(alertType string, config []byte) (account.Account, error) {
	return account.Account{}, nil
}
