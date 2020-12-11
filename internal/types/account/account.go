package account

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Account struct {
	// Unique identifier for  the account
	AccountID uuid.UUID `db:"account_id" json:"account_id,omitempty"`

	// API Key with full permission to the account
	RootAPIKey uuid.UUID `db:"root_api_key" json:"root_api_key,omitempty"`

	// The accounts alert type (slack, email, etc)
	AlertType string `db:"alert_type" json:"alert_type,omitempty"`

	// AlertConfig represents arbitrary json that will be marshalled
	// into the given alert type's configuration
	AlertConfig alertConfig `db:"alert_config" json:"alert_config,omitempty"`

	// Primary contact for the account
	AdminEmail string `db:"admin_email" json:"admin_email,omitempty"`
}

func (a Account) JSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a Account) String() (string, error) {
	x, err := a.JSON()
	return string(x), err
}
