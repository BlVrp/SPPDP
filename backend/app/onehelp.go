package app

import (
	fundraisestatuses "one-help/app/fundraise_statuses"
	"one-help/app/fundraises"
	"one-help/app/posts"
	"one-help/app/users"
)

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// Users provides access to users.DB.
	Users() users.DB

	// Posts provides access to posts.DB.
	Posts() posts.DB

	// FundraiseStatuses provides access to fundraise statuses DB.
	FundraiseStatuses() fundraisestatuses.DB

	// Fundraises provides access to fundraises DB.
	Fundraises() fundraises.DB

	// ExecuteMigrations applies migrations for the database.
	ExecuteMigrations(migrationsPath string, isUp bool) (err error)
}
