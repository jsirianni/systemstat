package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	b := "bob"
	p := 5555
	pgsql := Postgres{
		host:   b,
		port:   p,
		user:   b,
		pass:   b,
		dbname: b,
	}

	if err := pgsql.Validate(); err != nil {
		assert.Empty(t, err)
	}
	if err := pgsql.validateHost(); err != nil {
		assert.Empty(t, err)
	}
	if err := pgsql.validPort(); err != nil {
		assert.Empty(t, err)
	}
	if err := pgsql.validUser(); err != nil {
		assert.Empty(t, err)
	}
	if err := pgsql.validPassword(); err != nil {
		assert.Empty(t, err)
	}
	if err := pgsql.validDBName(); err != nil {
		assert.Empty(t, err)
	}
}
