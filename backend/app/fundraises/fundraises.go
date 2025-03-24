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
}

// IsEndDateSet returns true if end date is not null.
func (f *Fundraise) IsEndDateSet() bool {
	return f.EndDate != time.Time{}
}
