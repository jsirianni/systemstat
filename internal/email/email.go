package email

import (
    "github.com/badoux/checkmail"
)

func Validate(email string) error {
    return checkmail.ValidateFormat(email)
}
