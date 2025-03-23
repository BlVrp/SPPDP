package statuses

import (
	"context"
)

// DB exposes access to posts db.
//
// architecture: DB
type DB interface {
	// List returns all available posts.
	List(ctx context.Context) ([]string, error)
	// Create inserts post into the database.
	Create(ctx context.Context, post string) error
	// Delete post from the database.
	Delete(ctx context.Context, post string) error
}
