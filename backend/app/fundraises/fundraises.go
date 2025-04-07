package fundraises

import (
	"time"

	"github.com/google/uuid"
)

// Fundraise describes fundraise entity.
type Fundraise struct {
	ID           uuid.UUID
	OrganizerId  uuid.UUID
	Title        string
	Description  string
	TargetAmount float64
	StartDate    time.Time
	EndDate      time.Time
	Status       string
	ImageUrl     string
}

// IsEndDateSet returns true if end date is not null.
func (f *Fundraise) IsEndDateSet() bool {
	return f.EndDate != time.Time{}
}

// CreateParams defines needed params to create a new fundraise.
type CreateParams struct {
	OrganizerId  uuid.UUID
	Title        string
	Description  string
	TargetAmount float64
	EndDate      time.Time
	ImageUrl     string
}
