package database

import (
	"context"
	"database/sql"
	"errors"
	"one-help/app/donations"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrDonations indicates that there was an error in the database.
var ErrDonations = errs.Class("donations repository")

// donationsDB provides access to donations db.
//
// architecture: Database
type donationsDB struct {
	conn *sql.DB
}

// newDonationsDB is a constructor for base donationsDB.
func newDonationsDB(baseConn *sql.DB) donations.DB {
	return &donationsDB{
		conn: baseConn,
	}
}

// Create inserts donation into the database.
func (db *donationsDB) Create(ctx context.Context, donation donations.Donation) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrDonations.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO donations(donation_id, user_id, fundraise_id, amount, created_at)
              VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, donation.ID, donation.UserId, donation.FundraiseId, donation.Amount, donation.CreatedAt)
	return ErrDonations.Wrap(err)
}

// Get returns donation from the database by ID.
func (db *donationsDB) Get(ctx context.Context, id uuid.UUID) (donations.Donation, error) {
	var (
		donation donations.Donation
	)

	query := `SELECT donation_id, user_id, fundraise_id, amount, created_at
              FROM donations
              WHERE donation_id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&donation.ID,
		&donation.UserId,
		&donation.FundraiseId,
		&donation.Amount,
		&donation.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return donations.Donation{}, ErrDonations.Wrap(donations.ErrNoDonation)
		}
		return donation, ErrDonations.Wrap(err)
	}

	return donation, nil
}

// List returns all the donations.
func (db *donationsDB) List(ctx context.Context) ([]donations.Donation, error) {
	query := `SELECT donation_id, user_id, fundraise_id, amount, created_at
              FROM donations`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrDonations.Wrap(err)
	}
	defer rows.Close()

	var donationsList []donations.Donation
	for rows.Next() {
		var (
			donation donations.Donation
		)
		err := rows.Scan(
			&donation.ID,
			&donation.UserId,
			&donation.FundraiseId,
			&donation.Amount,
			&donation.CreatedAt,
		)
		if err != nil {
			return nil, ErrDonations.Wrap(err)
		}
		donationsList = append(donationsList, donation)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDonations.Wrap(err)
	}

	return donationsList, nil
}

// Update updates donation in database by id.
func (db *donationsDB) Update(ctx context.Context, donation donations.Donation) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrDonations.Wrap(err)
	}
	defer DeferCommitRollback(tx, &err)

	query := `UPDATE donations
	          SET user_id = $2, fundraise_id = $3, amount = $4, created_at = $5
	          WHERE donation_id = $1`

	_, err = tx.ExecContext(ctx, query,
		donation.ID,
		donation.UserId,
		donation.FundraiseId,
		donation.Amount,
		donation.CreatedAt,
	)
	if err != nil {
		return ErrDonations.Wrap(err)
	}

	return nil
}

// Delete removes donation from the database.
func (db *donationsDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM donations WHERE donation_id = $1`
	_, err := db.conn.ExecContext(ctx, query, id)
	return ErrFundraises.Wrap(err)
}
