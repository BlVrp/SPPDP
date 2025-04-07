package events

import (
	"time"

	"one-help/app/events"

	"github.com/google/uuid"
)

// CreateRequest defines parameters needed to create event.
type CreateRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	Format          string    `json:"format"`
	MaxParticipants int       `json:"maxParticipants"`
	MinimumDonation float64   `json:"minimumDonation"`
	Address         string    `json:"address"`
	Status          string    `json:"status"`
	FundraiseId     uuid.UUID `json:"fundraiseId"`
}

// EventView defines event view type.
type EventView struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	Format          string    `json:"format"`
	MaxParticipants int       `json:"maxParticipants"`
	MinimumDonation float64   `json:"minimumDonation"`
	Address         string    `json:"address"`
	Status          string    `json:"status"`
	FundraiseId     uuid.UUID `json:"fundraiseId"`
	CreatedAt       time.Time `json:"createdAt"`
}

// ToEventView converts event to view type.
func ToEventView(e *events.Event) EventView {
	return EventView{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		StartDate:       e.StartDate,
		EndDate:         e.EndDate,
		Format:          e.Format,
		MaxParticipants: e.MaxParticipants,
		MinimumDonation: e.MinimumDonation,
		Address:         e.Address,
		Status:          e.Status,
		FundraiseId:     e.FundraiseId,
		CreatedAt:       e.CreatedAt,
	}
}
