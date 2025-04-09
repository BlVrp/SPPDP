package events

import (
	"time"

	"github.com/google/uuid"
)

// Event describes event entity.
type Event struct {
	ID              uuid.UUID
	Title           string
	Description     string
	StartDate       time.Time
	EndDate         time.Time
	Format          string
	MaxParticipants int
	MinimumDonation float64
	Address         string
	Status          string
	FundraiseId     uuid.UUID
	CreatedAt       time.Time
	ImageUrl        string
}

// CreateParams defines needed params to create a new event.
type CreateParams struct {
	Title           string
	Description     string
	StartDate       time.Time
	EndDate         time.Time
	Format          string
	MaxParticipants int
	MinimumDonation float64
	Address         string
	FundraiseId     uuid.UUID
	ImageUrl        string
}
