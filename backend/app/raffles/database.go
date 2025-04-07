package raffles

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoRaffle indicates that raffle does not exist.
	ErrNoRaffle = errs.New("raffle does not exist")
)

// DB exposes access to raffles db.
//
// architecture: DB
type DB interface {
	// Create inserts raffle into the database.
	Create(ctx context.Context, raffle Raffle, gifts []Gift) error
	// Get raffle from the database.
	Get(ctx context.Context, id uuid.UUID) (Raffle, error)
	// List returns all available raffles.
	List(ctx context.Context, params ListParams) ([]Raffle, error)
	// ListGifts returns all raffle gifts.
	ListGifts(ctx context.Context, raffleID uuid.UUID) ([]Gift, error)
}

// ListParams defines params for list method.
type ListParams struct {
	FundraiseID *uuid.UUID
	Limit       int
	Page        int
}
