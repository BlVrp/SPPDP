package app

import (
	"one-help/app/users"
	"one-help/app/users/auth"
)

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// Users provides access to users.DB.
	Users() users.DB

	// Credentials provides access to auth.DB.
	Credentials() auth.DB

	// ExecuteMigrations applies migrations for the database.
	ExecuteMigrations(migrationsPath string, isUp bool) (err error)
}
