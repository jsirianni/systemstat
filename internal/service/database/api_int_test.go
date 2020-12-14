// +build integration

package database

import (
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var testIntServer Server

const testIntServerHost = "localhost"
const testIntServerPort = 19000

func init() {
	db, err := NewPostgres()
	if err != nil {
		panic(err)
	}

	testIntServer.Port.HTTP = testIntServerPort
	testIntServer.DB = db

	// run the server on a goroutine and then sleep two seconds
	// to make sure the server is running.
	go testIntServer.RunHTTP()
	if err := testConnect(); err != nil {
		panic(err)
	}
}

func testConnect() error {
	host := net.JoinHostPort(testIntServerHost, strconv.Itoa(testIntServerPort))
	attempts := 0
	for {
		_, err := net.DialTimeout("tcp", host, time.Second)
		if err == nil {
			break
		}
		if attempts > 3 {
			return errors.Wrap(err, "server failed to start on "+host)
		}
		time.Sleep(1)
	}
	return nil
}

func randomEmail() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + "@test.com"
}

func TestStatus(t *testing.T) {
	h, err := testIntServer.status()
	if err != nil {
		assert.Empty(t, err, "make sure the test server is running")
		return
	}
	assert.Equal(t, http.StatusOK, int(h.HttpStatus))
}

func TestGetAccount(t *testing.T) {
	id := "0234c572-15ec-4e67-9081-6a1c42a00090"
	email := "integration@test.com"

	a, err := testIntServer.getAccount(id)
	if err != nil {
		assert.Empty(t, err)
		return
	}
	assert.Equal(t, id, a.AccountId)
	assert.Equal(t, email, a.AdminEmail)
	assert.Equal(t, 200, int(a.HttpStatus))
}

func TestGetAccount404(t *testing.T) {
	id := "00000000-15ec-4e67-9081-6a1c42a00090"
	a, err := testIntServer.getAccount(id)
	if err != nil {
		assert.Equal(t, 404, int(a.HttpStatus))
		return
	}
	assert.Empty(t, err, "epected an error and 404 when getAccount is given a uuid that is not in the system")
}

// TODO: finish tests - createToken
