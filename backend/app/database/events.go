package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"one-help/app/events"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrEvents indicates that there was an error in the database.
var ErrEvents = errs.Class("events repository")

// eventsDB provides access to events db.
//
// architecture: Database
type eventsDB struct {
	conn *sql.DB
}

// newEventsDB is a constructor for base eventsDB.
func newEventsDB(baseConn *sql.DB) events.DB {
	return &eventsDB{
		conn: baseConn,
	}
}

// Create inserts event into the database.
func (db *eventsDB) Create(ctx context.Context, event events.Event) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrEvents.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO events(event_id, title, description, start_date, end_date, format, max_participants, minimum_donation, address, status, fundraise_id, created_at, form_url)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = tx.ExecContext(ctx, query, event.ID, event.Title, event.Description, event.StartDate, event.EndDate, event.Format, event.MaxParticipants, event.MinimumDonation, event.Address, event.Status, event.FundraiseId, event.CreatedAt, event.FormUrl)
	if err != nil {
		return ErrEvents.Wrap(err)
	}

	query = `INSERT INTO event_images(event_id, file_name) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, query, event.ID, event.ImageUrl)
	if err != nil {
		return ErrEvents.Wrap(err)
	}

	return nil
}

// Get returns event from the database by ID.
func (db *eventsDB) Get(ctx context.Context, id uuid.UUID) (events.Event, error) {
	var (
		event   events.Event
		endDate sql.NullTime
	)

	query := `SELECT events.event_id, title, description, start_date, end_date, format, max_participants, minimum_donation, address, status, fundraise_id, created_at, COALESCE(file_name, ''), form_url
              FROM events LEFT JOIN event_images ON events.event_id = event_images.event_id
              WHERE events.event_id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDate,
		&endDate,
		&event.Format,
		&event.MaxParticipants,
		&event.MinimumDonation,
		&event.Address,
		&event.Status,
		&event.FundraiseId,
		&event.CreatedAt,
		&event.ImageUrl,
		&event.FormUrl,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return events.Event{}, ErrEvents.Wrap(events.ErrNoEvents)
		}

		return event, ErrEvents.Wrap(err)
	}

	if endDate.Valid {
		event.EndDate = endDate.Time
	} else {
		event.EndDate = time.Time{}
	}

	return event, nil
}

// List returns all the events.
func (db *eventsDB) List(ctx context.Context, params events.ListParams) ([]events.Event, error) {
	var args = make([]any, 0, 3)
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Page == 0 {
		params.Page = 1
	}

	query := `SELECT events.event_id, title, description, start_date, end_date, format, max_participants, minimum_donation, address, status, fundraise_id, created_at, COALESCE(file_name, ''), form_url
              FROM events LEFT JOIN event_images ON events.event_id = event_images.event_id`

	{ // INFO: Paging.
		args = append(args, params.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
		args = append(args, (params.Page-1)*params.Limit)
		query += fmt.Sprintf(" OFFSET $%d ", len(args))
	}

	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, ErrEvents.Wrap(err)
	}
	defer func() { err = errs.Combine(err, rows.Close()) }()

	var eventsList []events.Event
	for rows.Next() {
		var (
			event   events.Event
			endDate sql.NullTime
		)
		err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&endDate,
			&event.Format,
			&event.MaxParticipants,
			&event.MinimumDonation,
			&event.Address,
			&event.Status,
			&event.FundraiseId,
			&event.CreatedAt,
			&event.ImageUrl,
			&event.FormUrl,
		)
		if err != nil {
			return nil, ErrEvents.Wrap(err)
		}

		if endDate.Valid {
			event.EndDate = endDate.Time
		} else {
			event.EndDate = time.Time{}
		}
		eventsList = append(eventsList, event)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrEvents.Wrap(err)
	}

	return eventsList, nil
}
