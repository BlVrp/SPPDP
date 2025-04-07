package events

import (
	"context"
	"time"

	"one-help/app/events/statuses"
	"one-help/internal/logger"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// Error wraps errors from users service that indicates about internal errors.
	Error = errs.Class("events service")
	// ParamsError wraps errors from users service that indicates about invalid or malformed parameters' data.
	ParamsError = errs.Class("users service: params")
)

// Service handles events related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger

	events DB
}

// NewService is a constructor for events service.
func NewService(logger logger.Logger, events DB) *Service {
	return &Service{
		logger: logger,
		events: events,
	}
}

// Create created new Event data in the system.
func (service *Service) Create(ctx context.Context, params CreateParams) (*Event, error) {
	switch {
	case params.Title == "":
		return nil, ParamsError.New("title is required")
	case params.Description == "":
		return nil, ParamsError.New("description is required")
	case params.StartDate.IsZero():
		return nil, ParamsError.New("start date is required")
	case params.Format == "":
		return nil, ParamsError.New("format is required")
	case params.MaxParticipants < 0:
		return nil, ParamsError.New("max participants must be positive or 0 for unlimited")
	case params.MinimumDonation < 0:
		return nil, ParamsError.New("minimum donation must be positive or 0 for no minimum")
	case params.Address == "":
		return nil, ParamsError.New("address is required")
	case params.FundraiseId == uuid.Nil:
		return nil, ParamsError.New("fundraise id is required")
	}

	event := &Event{
		ID:              uuid.New(),
		Title:           params.Title,
		Description:     params.Description,
		StartDate:       params.StartDate,
		EndDate:         params.EndDate,
		Format:          params.Format,
		MaxParticipants: params.MaxParticipants,
		MinimumDonation: params.MinimumDonation,
		Address:         params.Address,
		Status:          statuses.ActiveStatus,
		FundraiseId:     params.FundraiseId,
		CreatedAt:       time.Now().UTC(),
		ImageUrl:        params.ImageUrl,
	}

	err := service.events.Create(ctx, *event)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return event, nil
}

// Get returns Event by id.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (*Event, error) {
	event, err := service.events.Get(ctx, id)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &event, nil
}

// List returns list of events.
func (service *Service) List(ctx context.Context, limit, page int, creatorID *uuid.UUID) ([]Event, error) {
	switch {
	case limit <= 0:
		return nil, ParamsError.New("limit must be positive")
	case page <= 0:
		return nil, ParamsError.New("page must be positive")
	}

	list, err := service.events.List(ctx, ListParams{
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return list, nil
}
