package database

import (
	"context"
	"database/sql"
	fundraisestatuses "one-help/app/fundraise_statuses"

	"github.com/zeebo/errs"
)

// ErrFundraiseStatuses indicates that there was an error in the database.
var ErrFundraiseStatuses = errs.Class("fundraise statuses repository")

// fundraiseStatusesDB provides access to fundraise statuses db.
//
// architecture: Database
type fundraiseStatusesDB struct {
	conn *sql.DB
}

// newFundraiseStatusesDB is a constructor for base fundraiseStatusesDB.
func newFundraiseStatusesDB(baseConn *sql.DB) fundraisestatuses.DB {
	return &fundraiseStatusesDB{
		conn: baseConn,
	}
}

// Create inserts fundraise status into the database.
func (db *fundraiseStatusesDB) Create(ctx context.Context, fundraiseStatus string) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrFundraiseStatuses.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO fundraise_statuses(status)
              VALUES ($1)`
	_, err = tx.ExecContext(ctx, query, fundraiseStatus)
	if err != nil {
		return ErrFundraiseStatuses.Wrap(err)
	}

	return ErrFundraiseStatuses.Wrap(err)
}

// Delete removes fundraise status from the database.
func (db *fundraiseStatusesDB) Delete(ctx context.Context, fundraiseStatus string) error {
	query := `DELETE FROM fundraise_statuses WHERE status = $1`
	_, err := db.conn.ExecContext(ctx, query, fundraiseStatus)
	return ErrFundraiseStatuses.Wrap(err)
}

// List returns all available fundraise statuses.
func (db *fundraiseStatusesDB) List(ctx context.Context) ([]string, error) {
	query := `SELECT status FROM fundraise_statuses`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrFundraiseStatuses.Wrap(err)
	}

	defer rows.Close()

	var statuses []string
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			return nil, ErrFundraiseStatuses.Wrap(err)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
