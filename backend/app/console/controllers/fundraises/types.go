package fundraises

import (
	"time"

	"github.com/google/uuid"

	"one-help/app/fundraises"
	"one-help/app/users"
)

// CreateRequest defines request values for create endpoint.
type CreateRequest struct {
	OrganizerId  uuid.UUID `json:"organizerId"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	TargetAmount float64   `json:"targetAmount"`
	EndDate      time.Time `json:"endDate"`
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
	}
}

// LoginRequest defines request values for login endpoint.
type LoginRequest struct {
	Identifier string `json:"identifier"` // INFO: Email or phone number.
	Password   string `json:"password"`
}

// DonateRequest defines request values for donate endpoint.
type DonateRequest struct {
}

// UserPublicView defines user public view type.
type UserPublicView struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Website   string    `json:"website"`
	FileName  string    `json:"fileName"`
	City      string    `json:"city"`
}

// ToUserPublicView builds user public view.
func ToUserPublicView(user *users.User) *UserPublicView {
	return &UserPublicView{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Website:   user.Website,
		FileName:  user.FileName,
		City:      user.City,
	}
}
