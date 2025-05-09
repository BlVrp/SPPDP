package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeebo/errs"

	eventparticipants "one-help/app/events/participants"
)

// ErrEventParticipants indicates that there was an error in the database.
var ErrEventParticipants = errs.Class("event participants repository")

// eventParticipantsDB provides access to payments db.
//
// architecture: Database
type eventParticipantsDB struct {
	conn *sql.DB
}

// newEventParticipantsDB is a constructor for base eventParticipantsDB.
func newEventParticipantsDB(baseConn *sql.DB) eventparticipants.DB {
	return &eventParticipantsDB{
		conn: baseConn,
	}
}

// Create inserts event participant into the database.
func (db *eventParticipantsDB) Create(ctx context.Context, eventParticipant eventparticipants.EventParticipant) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrEventParticipants.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO event_participants(event_id, user_id)
              VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, query, eventParticipant.EventID, eventParticipant.UserID)
	if err != nil {
		return ErrEventParticipants.Wrap(err)
	}

	return ErrEventParticipants.Wrap(err)
}

// List returns all the event participants.
func (db *eventParticipantsDB) List(ctx context.Context, params eventparticipants.ListParams) ([]eventparticipants.EventParticipant, error) {
	var args = make([]any, 0, 4)

	var pagingEnabled = true

	if params.Limit == 0 {
		params.Limit = 20
		pagingEnabled = false
	}
	if params.Page == 0 {
		params.Page = 1
		pagingEnabled = false
	}

	query := `SELECT event_id, user_id
              FROM event_participants`

	var conditions []string

	if params.UserID != nil {
		args = append(args, *params.UserID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if params.EventID != nil {
		args = append(args, *params.EventID)
		conditions = append(conditions, fmt.Sprintf("event_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if pagingEnabled { // INFO: Paging.
		args = append(args, params.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
		args = append(args, (params.Page-1)*params.Limit)
		query += fmt.Sprintf(" OFFSET $%d ", len(args))
	}

	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, ErrEventParticipants.Wrap(err)
	}
	defer func() { err = errs.Combine(err, rows.Close()) }()

	var eventParticipantsList []eventparticipants.EventParticipant
	for rows.Next() {
		var (
			eventParticipant eventparticipants.EventParticipant
		)
		err = rows.Scan(
			&eventParticipant.UserID,
			&eventParticipant.EventID,
		)
		if err != nil {
			return nil, ErrEventParticipants.Wrap(err)
		}

		eventParticipantsList = append(eventParticipantsList, eventParticipant)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrEventParticipants.Wrap(err)
	}

	return eventParticipantsList, nil
}
