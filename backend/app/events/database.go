package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoEvents indicates that event does not exist.
	ErrNoEvents = errs.New("event does not exist")
)

// DB exposes access to fundraises db.
//
// architecture: DB
type DB interface {
	// Create inserts event into the database.
	Create(ctx context.Context, event Event) error
	// Get event from the database.
	Get(ctx context.Context, id uuid.UUID) (Event, error)
	// List returns all available events.
	List(ctx context.Context, params ListParams) ([]Event, error)
}

// ListParams defines params for list method.
type ListParams struct {
	Limit int
	Page  int
}
