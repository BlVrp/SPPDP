package app

import (
	"one-help/app/fundraises"
	fundraisestatuses "one-help/app/fundraises/statuses"
	"one-help/app/posts"
	"one-help/app/users"
	"one-help/app/users/credentials"
)

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// Users provides access to users.DB.
	Users() users.DB

	// Credentials provides access to credentials.DB.
	Credentials() credentials.DB

	// Posts provides access to posts.DB.
	Posts() posts.DB

	// FundraiseStatuses provides access to fundraise statuses DB.
	FundraiseStatuses() fundraisestatuses.DB

	// Fundraises provides access to fundraises DB.
	Fundraises() fundraises.DB

	// ExecuteMigrations applies migrations for the database.
	ExecuteMigrations(migrationsPath string, isUp bool) (err error)
}
