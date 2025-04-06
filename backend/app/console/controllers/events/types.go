package events

import (
	"one-help/app/events"
	"time"

	"github.com/google/uuid"
)

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
