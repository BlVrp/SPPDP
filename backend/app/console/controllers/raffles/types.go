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
	ImageUrl    string `json:"imageUrl"`
}

// RaffleView defines raffle view type.
type RaffleView struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	MinimumDonation float64    `json:"minimumDonation"`
	StartDate       time.Time  `json:"startDate"`
	EndDate         time.Time  `json:"endDate"`
	FundraiseID     uuid.UUID  `json:"fundraiseId"`
	Gifts           []GiftView `json:"gifts"`
}

// GiftView describes gift view.
type GiftView struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RaffleID    uuid.UUID `json:"raffleId"`
	UserID      uuid.UUID `json:"userId"`
	ImageUrl    string    `json:"imageUrl"`
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
			ImageUrl:    gift.ImageUrl,
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
