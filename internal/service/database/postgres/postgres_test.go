package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {
	p := Postgres{
		host:   "test",
		port:   9000,
		user:   "test",
		pass:   "test",
		dbname: "test",
	}

	actual := p.connectionString()
	expect := "host=test port=9000 user=test password=test dbname=test sslmode=disable"
	assert.Equal(t, expect, actual)
}
