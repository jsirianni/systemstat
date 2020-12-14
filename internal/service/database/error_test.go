package database

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsErrInvalidUUID(t *testing.T) {
	err := errors.New(`pq: invalid input syntax for type uuid: "not_a_uuid"`)
	if !isErrInvalidUUID(err) {
		t.Errorf("Expected isErrInvalidUUID to return true when given the error: " + err.Error())
		return
	}

	err = errors.New("random error")
	if isErrInvalidUUID(err) {
		t.Errorf("Expected isErrInvalidUUID to return false when given the error: " + err.Error())
	}
}

func TestisErrNoRows(t *testing.T) {
	var err = []error{}
	err = append(err, errors.New(`SELECT * FROM signup WHERE token = 'a98660b5-8304-470c-a931-7c84c4f314bd': sql: no rows in result set`))
	err = append(err, errors.New(`sql: no rows in result set`))

	for _, err := range err {
		if x := isErrNoRows(err); x != true {
			assert.Equal(t, true, x, err.Error())
		}
	}
}
