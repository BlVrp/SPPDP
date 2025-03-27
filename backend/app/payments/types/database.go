package types

import (
	"context"
)

// DB exposes access to payment types db.
//
// architecture: DB
type DB interface {
	// List returns all available payment types.
	List(ctx context.Context) ([]string, error)
	// Create inserts payment type into the database.
	Create(ctx context.Context, post string) error
	// Delete payment type from the database.
	Delete(ctx context.Context, post string) error
}
