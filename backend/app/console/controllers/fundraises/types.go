package fundraises

import (
	"time"

	"github.com/google/uuid"

	"one-help/app/fundraises"
)

// CreateRequest defines request values for create endpoint.
type CreateRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	TargetAmount float64   `json:"targetAmount"`
	EndDate      time.Time `json:"endDate"`
	ImageUrl     string    `json:"imageUrl"`
}

// FundraiseView defines fundraise view type.
type FundraiseView struct {
	ID           uuid.UUID `json:"id"`
	OrganizerId  uuid.UUID `json:"organizerId"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	TargetAmount float64   `json:"targetAmount"`
	FilledAmount float64   `json:"filledAmount"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate,omitempty"`
	Status       string    `json:"status"`
	ImageUrl     string    `json:"imageUrl"`
}

// ToFundraiseView builds fundraise view.
func ToFundraiseView(fundraise *fundraises.Fundraise, filled float64) *FundraiseView {
	return &FundraiseView{
		ID:           fundraise.ID,
		OrganizerId:  fundraise.OrganizerId,
		Title:        fundraise.Title,
		Description:  fundraise.Description,
		TargetAmount: fundraise.TargetAmount,
		FilledAmount: filled,
		StartDate:    fundraise.StartDate,
		EndDate:      fundraise.EndDate,
		Status:       fundraise.Status,
		ImageUrl:     fundraise.ImageUrl,
	}
}

// DonateRequest defines request values for donate endpoint.
type DonateRequest struct {
}
