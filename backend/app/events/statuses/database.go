package statuses

import (
	"context"
)

// DB exposes access to event statuses db.
//
// architecture: DB
type DB interface {
	// List returns all available event statuses.
	List(ctx context.Context) ([]string, error)
	// // Create inserts event status into the database.
	// Create(ctx context.Context, post string) error
	// // Delete event status from the database.
	// Delete(ctx context.Context, post string) error
}
