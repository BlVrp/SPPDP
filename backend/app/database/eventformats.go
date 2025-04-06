package database

import (
	"context"
	"database/sql"

	"github.com/zeebo/errs"

	eventformats "one-help/app/events/formats"
)

// ErrEventFormats indicates that there was an error in the database.
var ErrEventFormats = errs.Class("event formats repository")

// eventFormatsDB provides access to event formats db.
//
// architecture: Database
type eventFormatsDB struct {
	conn *sql.DB
}

// newEventFormatsDB is a constructor for base eventFormatsDB.
func newEventFormatsDB(baseConn *sql.DB) eventformats.DB {
	return &eventFormatsDB{
		conn: baseConn,
	}
}

// List returns all available event formats.
func (db *eventFormatsDB) List(ctx context.Context) ([]string, error) {
	query := `SELECT format FROM event_formats`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrEventFormats.Wrap(err)
	}

	defer func() { err = errs.Combine(err, rows.Close()) }()

	var statuses []string
	for rows.Next() {
		var status string
		if err = rows.Scan(&status); err != nil {
			return nil, ErrEventFormats.Wrap(err)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
