package raffles

import (
	"time"

	"github.com/google/uuid"
)

// Raffle describes raffle entity.
type Raffle struct {
	ID              uuid.UUID
	Title           string
	Description     string
	MinimumDonation float64
	StartDate       time.Time
	EndDate         time.Time
	FundraiseID     uuid.UUID
}

// CreateParams defines needed params to create a new raffle.
type CreateParams struct {
	Title           string
	Description     string
	MinimumDonation float64
	StartDate       time.Time
	EndDate         time.Time
	FundraiseID     uuid.UUID
	Gifts           []GiftCreateParams
}
