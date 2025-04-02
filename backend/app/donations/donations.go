package donations

import (
	"time"

	"github.com/google/uuid"
)

// Donation holds donate info.
type Donation struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	FundraiseId uuid.UUID
	Amount      float64
	CreatedAt   time.Time
}
