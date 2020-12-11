package email

/* We use a dedicated package for email handling because we may need
   to perform further evaluation on them in the future. */

import (
    "net/mail"
)

func Validate(email string) error {
    _, err := mail.ParseAddress(email)
    return err
}
