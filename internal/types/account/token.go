package account

import (
	"github.com/google/uuid"
)

type Token struct {
	// Signup token presented by users when creating an account
	Token uuid.UUID `db:"token" json:"token"`

	// false if unclaimed
	Claimed bool `db:"claimed" json:"claimed"`

	// Email or other unique identifier representing the user that claimed the token
	ClaimedBy string `db:"claimed_by" json:"claimed_by"`
}
