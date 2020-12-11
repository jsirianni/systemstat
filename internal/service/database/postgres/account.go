package postgres

import (
    "fmt"

    "github.com/jsirianni/systemstat/internal/types/account"

    "github.com/pkg/errors"
)

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
	err := p.queryAccount(q, &a)
	return a, err
}

func (p Postgres) AccountByEmail(email string) (account.Account, error) {
	a := account.Account{}

	if email == "" {
		return a, errors.New("email is a required parameter when reading an account")
	}

	q := fmt.Sprintf("SELECT * FROM account WHERE admin_email = '%s'", email)
	err := p.queryAccount(q, &a)
	return a, err
}

func (p Postgres) AccountConfigureAlert(alertType string, config []byte) (account.Account, error) {
	return account.Account{}, nil
}

func (p Postgres) queryAccount(q string, a *account.Account) error {
    return p.db.QueryRow(q).Scan(&a.AccountID, &a.RootAPIKey, &a.AlertType, &a.AlertConfig, &a.AdminEmail)
}
