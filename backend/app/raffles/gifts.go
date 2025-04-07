package raffles

import (
	"github.com/google/uuid"
)

// Gift describes gift entity.
type Gift struct {
	ID          uuid.UUID
	Title       string
	Description string
	RaffleID    uuid.UUID
	UserID      uuid.UUID
	ImageUrl    string
}

// GiftCreateParams defines needed parameters to create a gift.
type GiftCreateParams struct {
	Title       string
	Description string
	ImageUrl    string
}
