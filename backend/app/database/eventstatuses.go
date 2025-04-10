package database

import (
	"context"
	"database/sql"

	"github.com/zeebo/errs"

	eventstatuses "one-help/app/events/statuses"
)

// ErrEventStatuses indicates that there was an error in the database.
var ErrEventStatuses = errs.Class("event statuses repository")

// eventStatusesDB provides access to event statuses db.
//
// architecture: Database
type eventStatusesDB struct {
	conn *sql.DB
}

// neweventStatuesDB is a constructor for base eventStatuesDB.
func newEventStatusesDB(baseConn *sql.DB) eventstatuses.DB {
	return &eventStatusesDB{
		conn: baseConn,
	}
}

// List returns all available event statuses.
func (db *eventStatusesDB) List(ctx context.Context) ([]string, error) {
	query := `SELECT status FROM event_statuses`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrEventStatuses.Wrap(err)
	}

	defer func() { err = errs.Combine(err, rows.Close()) }()

	var statuses []string
	for rows.Next() {
		var status string
		if err = rows.Scan(&status); err != nil {
			return nil, ErrEventStatuses.Wrap(err)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
