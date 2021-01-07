// +build integration

package database

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/jsirianni/systemstat/internal/service/database/api"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStatusHTTP(t *testing.T) {
	uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/status"
	resp, err := http.Get(uri)
	if err != nil {
		assert.Empty(t, err)
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateTokenHTTP(t *testing.T) {
	uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/token/create"
	resp, err := http.Post(uri, "application/json", nil)
	if err != nil {
		assert.Empty(t, err)
		return
	}
	assert.Equal(t, 201, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		assert.Empty(t, err, "expected a json respsonse body from /v1/account/token/create")
		return
	}

	token := api.Token{}
	if err := json.Unmarshal(body, &token); err != nil {
		assert.Empty(t, err, "expected no errors when unmarshalling json response body into type Token")
		return
	}

	assert.NotEmpty(t, token.Token)
}

func TestCreateAccountHTTP(t *testing.T) {
	email := randomEmail()
	token, err := testIntServer.DB.CreateToken()
	if err != nil {
		assert.Empty(t, err, "expected CreateToken() to return a nil error")
		return
	}

	uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + token.Token + "/" + email
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

	a := api.Account{}
	if err := json.Unmarshal(body, &a); err != nil {
		assert.Empty(t, err, "expected no errors when unmarshalling json response body into type Account")
		return
	}

	assert.NotEmpty(t, a.AccountId)
	assert.NotEmpty(t, a.RootApiKey)
	assert.Equal(t, email, a.AdminEmail)

	// try a second time
	resp, err = http.Post(uri, "application/json", nil)
	if err != nil {
		assert.Empty(t, err)
		return
	}
	assert.Equal(t, 409, resp.StatusCode)
}

func TestGetAccountHTTP(t *testing.T) {
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

	a := api.Account{}
	if err := json.Unmarshal(body, &a); err != nil {
		assert.Empty(t, err, "expected no errors when unmarshalling json response body into type Account")
		return
	}

	assert.Equal(t, id, a.AccountId, string(body))
	assert.NotEmpty(t, a.RootApiKey, string(body))
	assert.Equal(t, email, a.AdminEmail, string(body))
}

func TestGetAccount404HTTP(t *testing.T) {
	uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + uuid.New().String()
	resp, err := http.Get(uri)
	if err != nil {
		assert.Empty(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	}
}

func TestGetAccount500HTTP(t *testing.T) {
	uri := "http://localhost:" + strconv.Itoa(testIntServerPort) + "/v1/account/" + "not_a_uuid"
	resp, err := http.Get(uri)
	if err != nil {
		assert.Empty(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	}
}
