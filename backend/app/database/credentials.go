package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeebo/errs"

	"one-help/app/users/credentials"
	"one-help/internal/postgres"
)

// ErrUserCredentials indicates that there was an error in the database.
var ErrUserCredentials = errs.Class("users repository")

// userCredentialsDB provides access to user credentials db.
//
// architecture: Database
type userCredentialsDB struct {
	conn *sql.DB
}

// newUserCredentialsDB is a constructor for base userCredentialsDB.
func newUserCredentialsDB(baseConn *sql.DB) credentials.DB {
	return &userCredentialsDB{
		conn: baseConn,
	}
}

// Create inserts user's credentials into the database.
func (db *userCredentialsDB) Create(ctx context.Context, creds credentials.Credentials) error {
	query := `INSERT INTO user_creds(user_id, phone_number, email, password_hash)
              VALUES ($1, $2, $3, $4)`
	_, err := db.conn.ExecContext(ctx, query, creds.UserID, creds.PhoneNumber, creds.Email, creds.PasswordHash)
	if err != nil {
		return parseUserCredentialsDBError(err)
	}

	return nil
}

// Get returns user's credentials from the database by user's ID.
func (db *userCredentialsDB) Get(ctx context.Context, key credentials.GetKey) (credentials.Credentials, error) {
	var creds credentials.Credentials
	query := `SELECT user_id, phone_number, email, password_hash
              FROM user_creds
              WHERE %s = $1`
	row := db.conn.QueryRowContext(ctx, fmt.Sprintf(query, key.FieldName()), key.Value())
	err := row.Scan(&creds.UserID, &creds.PhoneNumber, &creds.Email, &creds.PasswordHash)
	if err != nil {
		return creds, parseUserCredentialsDBError(err)
	}

	return creds, nil
}

// Update updates user's credentials in database by id.
func (db *userCredentialsDB) Update(ctx context.Context, creds credentials.Credentials) error {
	query := `UPDATE user_creds
              SET phone_number = $2, email = $3, password_hash = $4
              WHERE user_id = $1`
	reuslt, err := db.conn.ExecContext(ctx, query, creds.UserID, creds.PhoneNumber, creds.Email, creds.PasswordHash)
	if err != nil {
		return parseUserCredentialsDBError(err)
	}
	n, err := reuslt.RowsAffected()
	if err != nil {
		return ErrUserCredentials.Wrap(err)
	}
	if n == 0 {
		return ErrUserCredentials.Wrap(credentials.ErrNoUserCredentials)
	}

	return nil
}

// parseUserCredentialsDBError parses insert/update/... errors from db into typed package errors.
func parseUserCredentialsDBError(err error) error {
	switch {
	case postgres.IsConstraintError(err):
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "duplicate key value violates unique constraint \"user_creds_pkey\""):
			return ErrUserCredentials.Wrap(credentials.ErrUserAlreadyHasCredentials)
		case strings.Contains(errStr, `duplicate key value violates unique constraint "user_creds_email_key"`):
			return ErrUserCredentials.Wrap(credentials.ErrUserEmailTaken)
		case strings.Contains(errStr, `duplicate key value violates unique constraint "user_creds_phone_number_key"`):
			return ErrUserCredentials.Wrap(credentials.ErrUserPhoneNumberTaken)
		}
	case errs.Is(err, sql.ErrNoRows):
		return ErrUserCredentials.Wrap(credentials.ErrNoUserCredentials)
	}

	return ErrUserCredentials.Wrap(err)
}
