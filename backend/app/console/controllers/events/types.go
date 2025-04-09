package events

import (
	"time"

	"one-help/app/console/controllers/fundraises"
	"one-help/app/events"

	eventparticipants "one-help/app/events/participants"

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
	FundraiseId     uuid.UUID `json:"fundraiseId"`
	ImageUrl        string    `json:"imageUrl"`
}

// EventView defines event view type.
type EventView struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate,omitempty"`
	Format          string    `json:"format"`
	MaxParticipants int       `json:"maxParticipants"`
	MinimumDonation float64   `json:"minimumDonation"`
	Address         string    `json:"address"`
	Status          string    `json:"status"`
	FundraiseId     uuid.UUID `json:"fundraiseId"`
	CreatedAt       time.Time `json:"createdAt"`
	ImageUrl        string    `json:"imageUrl"`
}

// EventViewExtended defines event view type with additional data.
type EventViewExtended struct {
	EventView
	Fundraise fundraises.FundraiseView `json:"fundraise"`
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
		ImageUrl:        e.ImageUrl,
	}
}

// ToEventViewExtended converts event to extended view type.
func ToEventViewExtended(e *events.Event, fundraise fundraises.FundraiseView) EventViewExtended {
	return EventViewExtended{
		EventView: ToEventView(e),
		Fundraise: fundraise,
	}
}

type EventParticipantView struct {
	UserID  uuid.UUID `json:"userId"`
	EventID uuid.UUID `json:"eventId"`
}

func ToEventParticipantView(e *eventparticipants.EventParticipant) EventParticipantView {
	return EventParticipantView{
		UserID:  e.UserID,
		EventID: e.EventID,
	}
}
