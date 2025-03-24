package database

import (
	"context"
	"database/sql"
	"errors"
	"one-help/app/fundraises"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrFundraises indicates that there was an error in the database.
var ErrFundraises = errs.Class("fundraise repository")

// fundraisesDB provides access to fundraises db.
//
// architecture: Database
type fundraisesDB struct {
	conn *sql.DB
}

// newFundraisesDB is a constructor for base fundraisesDB.
func newFundraisesDB(baseConn *sql.DB) fundraises.DB {
	return &fundraisesDB{
		conn: baseConn,
	}
}

// Create inserts fundraise into the database.
func (db *fundraisesDB) Create(ctx context.Context, fundraise fundraises.Fundraise) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrFundraises.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO fundraises(fundraise_id, organizer_id, title, description, target_amount, start_date, end_date, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.ExecContext(ctx, query, fundraise.ID, fundraise.OrganizerId, fundraise.Title, fundraise.Description, fundraise.TargetAmount, fundraise.StartDate, fundraise.EndDate, fundraise.Status)
	return ErrFundraises.Wrap(err)
}

// Get returns fundraise from the database by ID.
func (db *fundraisesDB) Get(ctx context.Context, id uuid.UUID) (fundraises.Fundraise, error) {
	var (
		fundraise fundraises.Fundraise
		endDate   sql.NullTime
	)

	query := `SELECT id, organizer_id, title, description, target_amount, start_date, end_date, status
              FROM fundraises
              WHERE id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&fundraise.ID,
		&fundraise.OrganizerId,
		&fundraise.Title,
		&fundraise.Description,
		&fundraise.TargetAmount,
		&fundraise.StartDate,
		&endDate,
		&fundraise.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fundraises.Fundraise{}, ErrFundraises.Wrap(fundraises.ErrNoFundraise)
		}
		return fundraise, ErrFundraises.Wrap(err)
	}

	if endDate.Valid {
		fundraise.EndDate = endDate.Time
	} else {
		fundraise.EndDate = time.Time{}
	}

	return fundraise, nil
}

// List returns all the fundraises.
func (db *fundraisesDB) List(ctx context.Context) ([]fundraises.Fundraise, error) {
	query := `SELECT id, organizer_id, title, description, target_amount, start_date, end_date, status
              FROM fundraises`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrFundraises.Wrap(err)
	}
	defer rows.Close()

	var fundraisesList []fundraises.Fundraise
	for rows.Next() {
		var (
			fundraise fundraises.Fundraise
			endDate   sql.NullTime
		)
		err := rows.Scan(
			&fundraise.ID,
			&fundraise.OrganizerId,
			&fundraise.Title,
			&fundraise.Description,
			&fundraise.TargetAmount,
			&fundraise.StartDate,
			&endDate,
			&fundraise.Status,
		)
		if err != nil {
			return nil, ErrFundraises.Wrap(err)
		}

		if endDate.Valid {
			fundraise.EndDate = endDate.Time
		} else {
			fundraise.EndDate = time.Time{}
		}
		fundraisesList = append(fundraisesList, fundraise)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrFundraises.Wrap(err)
	}

	return fundraisesList, nil
}

func (db *fundraisesDB) Update(ctx context.Context, fundraise fundraises.Fundraise) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrFundraises.Wrap(err)
	}
	defer DeferCommitRollback(tx, &err)

	query := `UPDATE fundraises
	          SET organizer_id = $2, title = $3, description = $4, target_amount = $5, start_date = $6, end_date = $7, status = $8
	          WHERE id = $1`

	var endDate interface{}
	if fundraise.EndDate.IsZero() {
		endDate = nil
	} else {
		endDate = fundraise.EndDate
	}

	_, err = tx.ExecContext(ctx, query,
		fundraise.ID,
		fundraise.OrganizerId,
		fundraise.Title,
		fundraise.Description,
		fundraise.TargetAmount,
		fundraise.StartDate,
		endDate,
		fundraise.Status,
	)
	if err != nil {
		return ErrFundraises.Wrap(err)
	}

	return nil
}

// Delete removes fundraise from the database.
func (db *fundraisesDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM fundraises WHERE fundraise_id = $1`
	_, err := db.conn.ExecContext(ctx, query, id)
	return ErrFundraises.Wrap(err)
}
