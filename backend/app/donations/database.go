package donations

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrNoDonation indicates that donation does not exist.
var ErrNoDonation = errs.New("donation does not exist")

// DB exposes access to donations db.
//
// architecture: DB
type DB interface {
	// Create inserts donation into the database.
	Create(ctx context.Context, donation Donation) error
	// Get donation from the database.
	Get(ctx context.Context, id uuid.UUID) (Donation, error)
	// List returns all available donations.
	List(ctx context.Context, listParams ListParams) ([]Donation, error)
	// Update updates donation in database by id.
	Update(ctx context.Context, donation Donation) error
	// Delete donation from the database.
	Delete(ctx context.Context, id uuid.UUID) error
}
