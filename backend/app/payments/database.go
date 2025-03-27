package payments

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoPayment indicates that payment does not exist.
	ErrNoPayment = errs.New("payment does not exist")
)

// DB exposes access to payments db.
//
// architecture: DB
type DB interface {
	// Create inserts payment into the database.
	Create(ctx context.Context, payment Payment) error
	// Get payment from the database.
	Get(ctx context.Context, id uuid.UUID) (Payment, error)
	// List returns all available payments.
	List(ctx context.Context) ([]Payment, error)
	// Update updates payment in database by id.
	Update(ctx context.Context, payment Payment) error
	// Delete payment from the database.
	Delete(ctx context.Context, id uuid.UUID) error
}
