package email

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
    // based on https://gist.github.com/cjaoude/fd9910626629b53c4d25
    // but some addresses were wrong

    valid := []string{
        "email@example.com",
        "firstname.lastname@example.com",
        "email@subdomain.example.com",
        "firstname+lastname@example.com",
        "email@123.123.123.123",
        "1234567890@example.com",
        "email@example-one.com",
        "_______@example.com",
        "email@example.name",
        "email@example.museum",
        "email@example.co.jp",
        "firstname-lastname@example.com",
    }

    for _, email := range valid {
        err := Validate(email)
        assert.Empty(t, err, email)
    }

    invalid := []string{
        "plainaddress",
        "#@%^%#$@#$@#.com",
        "@example.com",
        "email.example.com",
        "email@example@example.com",
        ".email@example.com",
        "email.@example.com",
        "email..email@example.com",
        "email@example..com",
        "Abc..123@example.com",
    }

    for _, email := range invalid {
        err := Validate(email)
        assert.NotEmpty(t, err, email)
    }
}
