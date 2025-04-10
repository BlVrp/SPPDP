package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoUser indicates that user does not exist.
	ErrNoUser = errs.New("user does not exist")
	// ErrInvalidPassword indicates that invalid password was provided.
	ErrInvalidPassword = errs.New("invalid password")
)

// DB exposes access to users db.
//
// architecture: DB
type DB interface {
	// Create inserts user into the database.
	Create(ctx context.Context, user User) error
	// Get user from the database.
	Get(ctx context.Context, id uuid.UUID) (User, error)
	// Update updates user in database by id.
	Update(ctx context.Context, user User) error
	// Delete user from the database.
	Delete(ctx context.Context, id uuid.UUID) error
	// ListParticipants returns all raffle participants.
	ListRaffleParticipants(ctx context.Context, raffleID uuid.UUID) ([]UserWithContacts, error)
}
