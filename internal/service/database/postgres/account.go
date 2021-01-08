package postgres

import (
	"fmt"

	"github.com/jsirianni/systemstat/internal/service/database"

	"github.com/pkg/errors"
)

func (p Postgres) AccountCreate(email, token string) (database.Account, error) {
	if email == "" {
		return database.Account{}, errors.New("email is a required parameter when creating an account")
	}

	q := fmt.Sprintf("INSERT INTO account (admin_email) VALUES ('%s')", email)
	if _, err := p.db.Exec(q); err != nil {
		return database.Account{}, err
	}

	return p.AccountByEmail(email)
}

func (p Postgres) AccountByID(id string) (database.Account, error) {
	a := database.Account{}

	if id == "" {
		return a, errors.New("account_id is a required parameter when reading an account")
	}

	q := fmt.Sprintf("SELECT * FROM account WHERE account_id = '%s'", id)
	err := p.queryAccount(q, &a)
	return a, err
}

func (p Postgres) AccountByEmail(email string) (database.Account, error) {
	a := database.Account{}

	if email == "" {
		return a, errors.New("email is a required parameter when reading an account")
	}

	q := fmt.Sprintf("SELECT * FROM account WHERE admin_email = '%s'", email)
	err := p.queryAccount(q, &a)
	return a, err
}

func (p Postgres) AccountConfigureAlert(alertType string, config []byte) (database.Account, error) {
	return database.Account{}, nil
}

func (p Postgres) queryAccount(q string, a *database.Account) error {
	return p.db.QueryRow(q).Scan(&a.AccountId, &a.RootApiKey, &a.AlertType, &a.AlertConfig, &a.AdminEmail)
}
