package database

import (
	"context"
	"database/sql"

	"github.com/zeebo/errs"

	"one-help/app/payments/types"
)

// ErrPaymentTypes indicates that there was an error in the database.
var ErrPaymentTypes = errs.Class("payment types repository")

// paymentTypesDB provides access to payment types db.
//
// architecture: Database
type paymentTypesDB struct {
	conn *sql.DB
}

// newPaymentTypesDB is a constructor for base paymentTypesDB.
func newPaymentTypesDB(baseConn *sql.DB) types.DB {
	return &paymentTypesDB{
		conn: baseConn,
	}
}

// Create inserts payment type into the database.
func (db *paymentTypesDB) Create(ctx context.Context, paymentType string) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrPaymentTypes.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO payment_types(type)
              VALUES ($1)`
	_, err = tx.ExecContext(ctx, query, paymentType)
	if err != nil {
		return ErrPaymentTypes.Wrap(err)
	}

	return ErrPaymentTypes.Wrap(err)
}

// Delete removes payment type from the database.
func (db *paymentTypesDB) Delete(ctx context.Context, paymentType string) error {
	query := `DELETE FROM payment_types WHERE type = $1`
	_, err := db.conn.ExecContext(ctx, query, paymentType)
	return ErrPaymentTypes.Wrap(err)
}

// List returns all available payment types.
func (db *paymentTypesDB) List(ctx context.Context) ([]string, error) {
	query := `SELECT type FROM payment_types`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrPaymentTypes.Wrap(err)
	}

	defer func() { err = errs.Combine(err, rows.Close()) }()

	var types []string
	for rows.Next() {
		var paymentType string
		if err = rows.Scan(&paymentType); err != nil {
			return nil, ErrPaymentTypes.Wrap(err)
		}

		types = append(types, paymentType)
	}

	return types, nil
}
