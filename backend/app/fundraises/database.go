package fundraises

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoFundraise indicates that fundraise does not exist.
	ErrNoFundraise = errs.New("fundraise does not exist")
)

// DB exposes access to fundraises db.
//
// architecture: DB
type DB interface {
	// Create inserts fundraise into the database.
	Create(ctx context.Context, fundraise Fundraise) error
	// Get fundraise from the database.
	Get(ctx context.Context, id uuid.UUID) (Fundraise, error)
	// List returns all available fundraises.
	List(ctx context.Context, params ListParams) ([]Fundraise, error)
	// Update updates fundraise in database by id.
	Update(ctx context.Context, fundraise Fundraise) error
	// Delete fundraise from the database.
	Delete(ctx context.Context, id uuid.UUID) error
	// GetFilled returns collected funds on the provided fundraise.
	GetFilled(ctx context.Context, id uuid.UUID) (float64, error)
}

// ListParams defines params for list method.
type ListParams struct {
	OrganizerID *uuid.UUID
	Limit       int
	Page        int
}
