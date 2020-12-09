// +build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	pgsql := Postgres{}
	if err := pgsql.configure(); err != nil {
		assert.Empty(t, err)
		return
	}
}
