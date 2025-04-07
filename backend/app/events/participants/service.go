package participants

import (
	"context"
	"one-help/internal/logger"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// Error wraps errors from event participants service that indicates about internal errors.
	Error = errs.Class("event participants service")
	// ParamsError wraps errors from event participants service that indicates about invalid or malformed parameters' data.
	ParamsError = errs.Class("event participants service: params")
)

// Service handles event participants related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger

	eventParticipants DB
}

// NewService is a constructor for event participant service.
func NewService(logger logger.Logger, eventParticipants DB) *Service {
	return &Service{
		logger:            logger,
		eventParticipants: eventParticipants,
	}
}

// Create creates new Event Participant data in the system.
func (service *Service) Create(ctx context.Context, ep EventParticipant) (*EventParticipant, error) {

	err := service.eventParticipants.Create(ctx, ep)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &ep, nil
}

// List returns list of event participants.
func (service *Service) List(ctx context.Context, limit, page int, userID *uuid.UUID) ([]EventParticipant, error) {
	switch {
	case limit <= 0:
		return nil, ParamsError.New("limit must be positive")
	case page <= 0:
		return nil, ParamsError.New("page must be positive")
	}

	list, err := service.eventParticipants.List(ctx, ListParams{
		UserID: userID,
		Limit:  limit,
		Page:   page,
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return list, nil
}
