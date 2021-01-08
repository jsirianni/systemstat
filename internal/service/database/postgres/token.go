package postgres

import (
	"fmt"

	"github.com/jsirianni/systemstat/internal/service/database"

	"github.com/pkg/errors"
)

func (p Postgres) ClaimToken(email, token string) (database.Token, error) {
	// check if token exists and is claimed
	t, err := p.GetToken(token)
	if err != nil {
		return t, err
	}
	if t.Claimed {
		return t, errors.New("token " + token + " is already claimed")
	}

	// claim the token
	q := fmt.Sprintf("UPDATE signup SET claimed = true, claimed_by = '%s' WHERE token = '%s' AND claimed = false", email, token)
	if _, err := p.db.Exec(q); err != nil {
		return t, err
	}

	// validate that the token has been claimed
	t, err = p.GetToken(token)
	if err != nil {
		return t, err
	}
	if !t.Claimed {
		return t, errors.New("token should be claimed, got false")
	}
	if t.ClaimedBy != email {
		return t, errors.New("token should be claimed by " + email + " but got " + t.ClaimedBy)
	}

	return t, err
}

func (p Postgres) GetToken(token string) (database.Token, error) {
	t := database.Token{}

	if token == "" {
		return t, errors.New("token is a required parameter when retrieving a token")
	}

	q := fmt.Sprintf("SELECT * FROM signup WHERE token = '%s'", token)
	err := p.queryToken(q, &t)
	return t, err
}

func (p Postgres) CreateToken() (database.Token, error) {
	t := database.Token{}
	q := fmt.Sprintf("INSERT INTO signup DEFAULT VALUES RETURNING token, claimed, claimed_by")
	err := p.queryToken(q, &t)
	return t, err
}

func (p Postgres) queryToken(q string, t *database.Token) error {
	err := p.db.QueryRow(q).Scan(&t.Token, &t.Claimed, &t.ClaimedBy)
	return errors.Wrap(err, q)
}
