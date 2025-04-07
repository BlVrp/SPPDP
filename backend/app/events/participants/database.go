package participants

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoEventParticipants indicates that event does not exist.
	ErrNoEventsParticipants = errs.New("event participant does not exist")
)

// DB exposes access to event participants db.
//
// architecture: DB
type DB interface {
	// Create inserts event participant into the database.
	Create(ctx context.Context, eventParticipant EventParticipant) error
	// List returns all available event participants.
	List(ctx context.Context, params ListParams) ([]EventParticipant, error)
}

// ListParams defines params for list method.
type ListParams struct {
	UserID *uuid.UUID
	Limit  int
	Page   int
}
