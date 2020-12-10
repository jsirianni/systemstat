package postgres

import (
	"fmt"

	"github.com/jsirianni/systemstat/internal/types/account"

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
	return account.Account{}, nil
}

func (p Postgres) AccountByEmail(email string) (account.Account, error) {
	q := "SELECT * FROM account WHERE admin_email = 'slack@test.com'"
	a := account.Account{}
	err := p.db.QueryRow(q).Scan(&a.AccountID, &a.RootAPIKey, &a.AlertType, &a.AlertConfig, &a.AdminEmail)
	return a, err
}

func (p Postgres) AccountConfigureAlert(alertType string, config []byte) (account.Account, error) {
	return account.Account{}, nil
}
