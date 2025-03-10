package app

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// ExecuteMigrations applies migrations for the database.
	ExecuteMigrations(migrationsPath string, isUp bool) (err error)
}
