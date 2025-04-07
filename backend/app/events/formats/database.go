package formats

import "context"

// DB exposes access to event formats db.
//
// architecture: DB
type DB interface {
	// List returns all available event formats.
	List(ctx context.Context) ([]string, error)
	// // Create inserts fundraise status into the database.
	// Create(ctx context.Context, post string) error
	// // Delete fundraise status from the database.
	// Delete(ctx context.Context, post string) error
}
