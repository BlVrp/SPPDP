package raffles

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"one-help/internal/logger"
)

var (
	// Error wraps errors from raffles service that indicates about internal errors.
	Error = errs.Class("raffles service")
	// ParamsError wraps errors from raffles service that indicates about invalid or malformed parameters' data.
	ParamsError = errs.Class("raffles service: params")
)

// Service handles raffles related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger

	raffles DB
}

// NewService is a constructor for raffles service.
func NewService(logger logger.Logger, raffles DB) *Service {
	return &Service{
		logger:  logger,
		raffles: raffles,
	}
}

// Create creates new Raffle data in the system.
func (service *Service) Create(ctx context.Context, params CreateParams) (*Raffle, []Gift, error) {
	switch {
	case params.Title == "":
		return nil, nil, ParamsError.New("title is required")
	case params.Description == "":
		return nil, nil, ParamsError.New("description is required")
	case params.MinimumDonation < 0.:
		return nil, nil, ParamsError.New("minimum amount must be positive on zero")
	}

	raffle := &Raffle{
		ID:              uuid.New(),
		Title:           params.Title,
		Description:     params.Description,
		MinimumDonation: params.MinimumDonation,
		StartDate:       params.StartDate,
		EndDate:         params.EndDate,
		FundraiseID:     params.FundraiseID,
	}

	gifts := make([]Gift, len(params.Gifts))
	for i, gift := range params.Gifts {
		gifts[i] = Gift{
			ID:          uuid.New(),
			Title:       gift.Title,
			Description: gift.Description,
			RaffleID:    raffle.ID,
			UserID:      uuid.Nil,
			ImageUrl:    gift.ImageUrl,
		}
	}

	err := service.raffles.Create(ctx, *raffle, gifts)
	if err != nil {
		return nil, nil, Error.Wrap(err)
	}

	return raffle, gifts, nil
}

// Get returns Raffle by id.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (*Raffle, error) {
	raffle, err := service.raffles.Get(ctx, id)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &raffle, nil
}

// List returns list of raffles.
func (service *Service) List(ctx context.Context, limit, page int, fundraiseID *uuid.UUID) ([]Raffle, error) {
	switch {
	case limit <= 0:
		return nil, ParamsError.New("limit must be positive")
	case page <= 0:
		return nil, ParamsError.New("page must be positive")
	}

	list, err := service.raffles.List(ctx, ListParams{
		FundraiseID: fundraiseID,
		Limit:       limit,
		Page:        page,
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return list, nil
}

// ListGifts returns list of raffle gifts.
func (service *Service) ListGifts(ctx context.Context, raffleID uuid.UUID) ([]Gift, error) {
	list, err := service.raffles.ListGifts(ctx, raffleID)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return list, nil
}
