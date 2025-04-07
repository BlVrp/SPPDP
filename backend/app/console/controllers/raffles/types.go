package raffles

import (
	"time"

	"github.com/google/uuid"

	"one-help/app/raffles"
)

// CreateRequest defines request values for create endpoint.
type CreateRequest struct {
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	MinimumDonation float64             `json:"minimumDonation"`
	StartDate       time.Time           `json:"startDate"`
	EndDate         time.Time           `json:"endDate"`
	FundraiseID     uuid.UUID           `json:"fundraiseId"`
	Gifts           []CreateGiftRequest `json:"gifts"`
}

// CreateGiftRequest defines request values needed to create gift.
type CreateGiftRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"imageLink"`
}

// RaffleView defines raffle view type.
type RaffleView struct {
	ID              uuid.UUID
	Title           string
	Description     string
	MinimumDonation float64
	StartDate       time.Time
	EndDate         time.Time
	FundraiseID     uuid.UUID
	Gifts           []GiftView
}

// GiftView describes gift view.
type GiftView struct {
	ID          uuid.UUID
	Title       string
	Description string
	RaffleID    uuid.UUID
	UserID      uuid.UUID
	ImageLink   string
}

// ToRaffleView builds raffle view.
func ToRaffleView(raffle *raffles.Raffle, gifts []raffles.Gift) *RaffleView {
	var giftsView = make([]GiftView, len(gifts))
	for i, gift := range gifts {
		giftsView[i] = GiftView{
			ID:          gift.ID,
			Title:       gift.Title,
			Description: gift.Description,
			RaffleID:    gift.RaffleID,
			UserID:      gift.UserID,
			ImageLink:   gift.ImageLink,
		}
	}

	return &RaffleView{
		ID:              raffle.ID,
		Title:           raffle.Title,
		Description:     raffle.Description,
		MinimumDonation: raffle.MinimumDonation,
		StartDate:       raffle.StartDate,
		EndDate:         raffle.EndDate,
		FundraiseID:     raffle.FundraiseID,
		Gifts:           giftsView,
	}
}
