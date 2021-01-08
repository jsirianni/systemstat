package database

type Database interface {
	// Validate the database configuration
	Validate() error

	// Test connection to the database
	TestConnection() error

	// Create an account
	AccountCreate(email, token string) (Account, error)

	// Retrieve an account by account_id
	AccountByID(id string) (Account, error)

	// Retrieve an account by email
	AccountByEmail(email string) (Account, error)

	// Configure an accounts alert type
	AccountConfigureAlert(alertType string, config []byte) (Account, error)

	// claim a sign up token
	ClaimToken(email, token string) (Token, error)

	// get an existing token
	GetToken(token string) (Token, error)

	// create a signup token
	CreateToken() (Token, error)
}
