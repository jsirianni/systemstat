// +build integration

package database

import (
    "strconv"
    "net"
    "net/http"
    "testing"
    "time"
    "encoding/json"
    "io/ioutil"

    "github.com/jsirianni/systemstat/internal/types/account"

    "github.com/stretchr/testify/assert"
    "github.com/pkg/errors"
)

var testIntServer Server
const testIntServerHost = "localhost"
const testIntServerPort = 19000

func init() {
    db, err := NewPostgres()
    if err != nil {
        panic(err)
    }

    testIntServer.Port = testIntServerPort
    testIntServer.DB   = db

    // run the server on a goroutine and then sleep two seconds
    // to make sure the server is running.
    go testIntServer.Run()
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
            return errors.Wrap(err, "server failed to start on " + host)
        }
        time.Sleep(1)
    }
    return nil
}

func randomEmail() string {
    return strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + "@test.com"
}

func TestStatus(t *testing.T) {
    uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/status"
    resp, err := http.Get(uri)
    if err != nil {
        assert.Empty(t, err)
        return
    }
    assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateAccount(t *testing.T) {
    email := randomEmail()
    uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + email
    resp, err := http.Post(uri, "application/json", nil)
    if err != nil {
        assert.Empty(t, err)
        return
    }
    assert.Equal(t, 201, resp.StatusCode)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        assert.Empty(t, err, "expected a json respsonse body from /v1/account/{account}")
        return
    }

    a := account.Account{}
    if err := json.Unmarshal(body, &a); err != nil {
        assert.Empty(t, err, "expected no errors when unmarshalling json response body into type Account")
        return
    }

    assert.NotEmpty(t, a.AccountID)
    assert.NotEmpty(t, a.RootAPIKey)
    assert.Equal(t, email, a.AdminEmail)

    // try a second time
    resp, err = http.Post(uri, "application/json", nil)
    if err != nil {
        assert.Empty(t, err)
        return
    }
    assert.Equal(t, 409, resp.StatusCode)
}

func TestGetAccount(t *testing.T) {
    // id is from scripts/postgres/test_data.sql
    id := "0234c572-15ec-4e67-9081-6a1c42a00090"
    email := "integration@test.com"

    uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + id
    resp, err := http.Get(uri)
    if err != nil {
        assert.Empty(t, err)
        return
    }
    assert.Equal(t, 200, resp.StatusCode)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        assert.Empty(t, err, "expected a json respsonse body from /v1/account/{account}")
        return
    }

    a := account.Account{}
    if err := json.Unmarshal(body, &a); err != nil {
        assert.Empty(t, err, "expected no errors when unmarshalling json response body into type Account")
        return
    }

    assert.NotEqual(t, id, a.AccountID)
    assert.NotEmpty(t, a.RootAPIKey)
    assert.Equal(t, email, a.AdminEmail)
}

func TestGetAccount404(t *testing.T) {
    uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + randomEmail()
    resp, err := http.Get(uri)
    if err != nil {
        assert.Empty(t, err)
        return
    }
    assert.Equal(t, 404, resp.StatusCode)
}
