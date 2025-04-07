package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // using golang migrate source.
	_ "github.com/lib/pq"
	"github.com/zeebo/errs"

	"one-help/app"
	"one-help/app/donations"
	"one-help/app/events"
	"one-help/app/fundraises"
	fundraisestatuses "one-help/app/fundraises/statuses"
	"one-help/app/payments"
	paymenttypes "one-help/app/payments/types"
	"one-help/app/posts"
	"one-help/app/users"
	"one-help/app/users/credentials"
)

var logger = log.Default()

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

// Users provides access to users.DB.
func (db *database) Users() users.DB {
	return newUsersDB(db.conn)
}

// Credentials provides access to credentials.DB.
func (db *database) Credentials() credentials.DB {
	return newUserCredentialsDB(db.conn)
}

// Posts provides access to posts.DB.
func (db *database) Posts() posts.DB {
	return newPostsDB(db.conn)
}

// FundraiseStatuses provides access to fundraise statuses DB.
func (db *database) FundraiseStatuses() fundraisestatuses.DB {
	return newFundraiseStatusesDB(db.conn)
}

// Fundraises provides access to fundraises DB.
func (db *database) Fundraises() fundraises.DB {
	return newFundraisesDB(db.conn)
}

// Events provides access to events DB.
func (db *database) Events() events.DB {
	return newEventsDB(db.conn)
}

// Donations provides access to donations DB.
func (db *database) Donations() donations.DB {
	return newDonationsDB(db.conn)
}

// PaymentTypes provides access to payment types DB.
func (db *database) PaymentTypes() paymenttypes.DB {
	return newPaymentTypesDB(db.conn)
}

// Payments provides access to payments DB.
func (db *database) Payments() payments.DB {
	return newPaymentsDB(db.conn)
}

// Close closes underlying db connection.
func (db *database) Close() error {
	return Error.Wrap(db.conn.Close())
}

// DeferCommitRollback provides helping functionality to commit or rollback transaction,
// based on the error state, with provided error updated.
// NOTE: in case the DeferCommitRollback called as a deffer statement,
// Rollback or Commit should not be called by themselves.
// NOTE: if DeferCommitRollback is called from function that panics, the panic value will be recovered,
// database transaction will be closed with rollback, and panic value will be printed to stdout.
func DeferCommitRollback(tx *sql.Tx, err *error) {
	if panicVal := recover(); panicVal != nil {
		logger.Println("[DeferCommitRollback] called from function that panics -> database transaction will be closed with rollback")
		logger.Println(panicVal) // INFO: Notify about panic message.
		*err = errs.Combine(*err, fmt.Errorf("panic detected"))
	}

	var innerErr error
	if *err != nil {
		innerErr = tx.Rollback()
		if innerErr != nil {
			logger.Println("failed to rollback transaction", innerErr)
		}
	} else {
		innerErr = tx.Commit()
		if innerErr != nil {
			logger.Println("failed to commit transaction", innerErr)
		}
	}
	*err = errs.Combine(*err, innerErr)
}
