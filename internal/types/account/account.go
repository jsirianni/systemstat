package account

import (
    "github.com/google/uuid"
)

type Account struct {
    // Unique identifier for  the account
    AccountID  uuid.UUID `db:"account_id" json:"account_id"`

    // API Key with full permission to the account
    RootAPIKey uuid.UUID `db:"root_api_key" json:"root_api_key"`

    // The accounts alert type (slack, email, etc)
    AlertType  string    `db:"alert_type" json:"alert_type"`

    // AlertConfig represents arbitrary json that will be marshalled
    // into the given alert type's configuration
    AlertConfig map[string]interface{} `db:"alert_config" json:"alert_config"`

    // Primary contact for the account
    AdminEmail string `db:"admin_email" json:"admin_email"`
}
