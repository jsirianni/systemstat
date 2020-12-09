// +build integration

package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// global integration test Postgres instance
var testPG Postgres

func init() {
	var err error
	testPG, err = New()
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		assert.Empty(t, err)
	}
}

func TestTestConnection(t *testing.T) {
	if err := testPG.TestConnection(); err != nil {
		assert.Empty(t, err)
	}
}
