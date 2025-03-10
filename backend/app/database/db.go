package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // using golang migrate source.
	_ "github.com/lib/pq"
	"github.com/zeebo/errs"

	"one-help/app"
)

// Config contains configurable values for one-help db.
type Config struct {
	Database       string `env:"DATABASE"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

// ensures that database implements app.DB.
var _ app.DB = (*database)(nil)

var (
	// Error is the default one-help db error class.
	Error = errs.Class("master database")
)

// database combines access to different database tables with a record
// of the db driver, db implementation, and db source URL.
//
// architecture: Master Database
type database struct {
	conn *sql.DB
}

// New returns app.DB postgresql implementation.
func New(databaseURL string) (app.DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &database{conn: conn}, nil
}

// ExecuteMigrations executes migrations by path in database.
func (db *database) ExecuteMigrations(migrationsPath string, isUp bool) (err error) {
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return Error.Wrap(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		return Error.Wrap(err)
	}

	if isUp {
		return Error.Wrap(m.Up())
	}

	return Error.Wrap(m.Down())
}

// Close closes underlying db connection.
func (db *database) Close() error {
	return Error.Wrap(db.conn.Close())
}
