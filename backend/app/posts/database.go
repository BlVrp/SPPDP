package posts

import (
	"context"
)

// DB exposes access to posts db.
//
// architecture: DB
type DB interface {
	// List returns all available posts.
	List(ctx context.Context) ([]string, error)
}
