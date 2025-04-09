package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"one-help/app/raffles"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrRaffles indicates that there was an error in the database.
var ErrRaffles = errs.Class("raffles repository")

// rafflesDB provides access to raffles db.
//
// architecture: Database
type rafflesDB struct {
	conn *sql.DB
}

// newRafflesDB is a constructor for base rafflesDB.
func newRafflesDB(baseConn *sql.DB) raffles.DB {
	return &rafflesDB{
		conn: baseConn,
	}
}

// Create inserts raffle into the database.
func (db *rafflesDB) Create(ctx context.Context, raffle raffles.Raffle, gifts []raffles.Gift) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrRaffles.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO raffles(raffle_id, title, description, minimum_donation, start_date, end_date, fundraise_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, query, raffle.ID, raffle.Title, raffle.Description, raffle.MinimumDonation, raffle.StartDate, raffle.EndDate, raffle.FundraiseID)
	if err != nil {
		return ErrRaffles.Wrap(err)
	}

	giftQuery := `INSERT INTO gifts(gift_id, title, description, raffle_id, user_id) VALUES ($1, $2, $3, $4, $5)`
	giftStmt, err := tx.PrepareContext(ctx, giftQuery)
	if err != nil {
		return ErrRaffles.Wrap(err)
	}

	defer func() { err = errs.Combine(err, giftStmt.Close()) }()

	giftImageQuery := `INSERT INTO gift_images(gift_id, file_name) VALUES ($1, $2)`
	giftImageStmt, err := tx.PrepareContext(ctx, giftImageQuery)
	if err != nil {
		return ErrRaffles.Wrap(err)
	}

	defer func() { err = errs.Combine(err, giftImageStmt.Close()) }()

	for _, gift := range gifts {
		_, err = giftStmt.ExecContext(ctx, gift.ID, gift.Title, gift.Description, gift.RaffleID, nil)
		if err != nil {
			return ErrRaffles.Wrap(err)
		}

		_, err = giftImageStmt.ExecContext(ctx, gift.ID, gift.ImageUrl)
		if err != nil {
			return ErrRaffles.Wrap(err)
		}
	}

	return ErrRaffles.Wrap(err)
}

// Get returns raffle from the database by ID.
func (db *rafflesDB) Get(ctx context.Context, id uuid.UUID) (raffles.Raffle, error) {
	var raffle raffles.Raffle
	query := `SELECT raffle_id, title, description, minimum_donation, start_date, end_date, fundraise_id
              FROM raffles
              WHERE raffle_id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&raffle.ID,
		&raffle.Title,
		&raffle.Description,
		&raffle.MinimumDonation,
		&raffle.StartDate,
		&raffle.EndDate,
		&raffle.FundraiseID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return raffles.Raffle{}, ErrRaffles.Wrap(raffles.ErrNoRaffle)
		}

		return raffle, ErrRaffles.Wrap(err)
	}

	return raffle, nil
}

// List returns all the raffles.
func (db *rafflesDB) List(ctx context.Context, params raffles.ListParams) ([]raffles.Raffle, error) {
	var args = make([]any, 0, 3)
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Page == 0 {
		params.Page = 1
	}

	query := `SELECT raffle_id, title, description, minimum_donation, start_date, end_date, fundraise_id
              FROM raffles`

	if params.FundraiseID != nil {
		args = append(args, *params.FundraiseID)
		query += fmt.Sprintf(" WHERE fundraize_id = $%d ", len(args))
	}

	{ // INFO: Paging.
		args = append(args, params.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
		args = append(args, (params.Page-1)*params.Limit)
		query += fmt.Sprintf(" OFFSET $%d ", len(args))
	}

	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, ErrRaffles.Wrap(err)
	}
	defer func() { err = errs.Combine(err, rows.Close()) }()

	var rafflesList []raffles.Raffle
	for rows.Next() {
		var raffle raffles.Raffle
		err = rows.Scan(
			&raffle.ID,
			&raffle.Title,
			&raffle.Description,
			&raffle.MinimumDonation,
			&raffle.StartDate,
			&raffle.EndDate,
			&raffle.FundraiseID,
		)
		if err != nil {
			return nil, ErrRaffles.Wrap(err)
		}

		rafflesList = append(rafflesList, raffle)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrRaffles.Wrap(err)
	}

	return rafflesList, nil
}

// ListGifts returns all raffle gifts.
func (db *rafflesDB) ListGifts(ctx context.Context, raffleID uuid.UUID) ([]raffles.Gift, error) {
	query := `SELECT gifts.gift_id, title, description, raffle_id, user_id, COALESCE(file_name, '')
              FROM gifts LEFT JOIN gift_images ON gifts.gift_id = gift_images.gift_id
              WHERE raffle_id = $1`

	rows, err := db.conn.QueryContext(ctx, query, raffleID)
	if err != nil {
		return nil, ErrRaffles.Wrap(err)
	}
	defer func() { err = errs.Combine(err, rows.Close()) }()

	var giftsList []raffles.Gift
	for rows.Next() {
		var gift raffles.Gift
		err = rows.Scan(
			&gift.ID,
			&gift.Title,
			&gift.Description,
			&gift.RaffleID,
			&gift.UserID,
			&gift.ImageUrl,
		)
		if err != nil {
			return nil, ErrRaffles.Wrap(err)
		}

		giftsList = append(giftsList, gift)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrRaffles.Wrap(err)
	}

	return giftsList, nil
}
