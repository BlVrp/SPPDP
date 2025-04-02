package fundraises

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"one-help/app/fundraises/statuses"
	"one-help/internal/logger"
)

var (
	// Error wraps errors from users service that indicates about internal errors.
	Error = errs.Class("fundraises service")
	// ParamsError wraps errors from users service that indicates about invalid or malformed parameters' data.
	ParamsError = errs.Class("users service: params")
)

// Service handles fundraises related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger

	fundraises DB
}

// NewService is a constructor for fundraises service.
func NewService(logger logger.Logger, fundraises DB) *Service {
	return &Service{
		logger:     logger,
		fundraises: fundraises,
	}
}

// Create created new Fundraise data in the system.
func (service *Service) Create(ctx context.Context, params CreateParams) (*Fundraise, error) {
	switch {
	case params.Title == "":
		return nil, ParamsError.New("title is required")
	case params.Description == "":
		return nil, ParamsError.New("description is required")
	case params.TargetAmount <= 0.:
		return nil, ParamsError.New("target amount must be positive")
	}

	fundraise := &Fundraise{
		ID:           uuid.New(),
		OrganizerId:  params.OrganizerId,
		Title:        params.Title,
		Description:  params.Description,
		TargetAmount: params.TargetAmount,
		StartDate:    time.Now().UTC(),
		EndDate:      params.EndDate,
		Status:       statuses.ActiveStatus,
	}

	err := service.fundraises.Create(ctx, *fundraise)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return fundraise, nil
}

// Get returns Fundraise be id.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (*Fundraise, error) {
	fundraise, err := service.fundraises.Get(ctx, id)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &fundraise, nil
}

// List returns list of fundraises.
func (service *Service) List(ctx context.Context, limit, page int, creatorID *uuid.UUID) ([]Fundraise, error) {
	switch {
	case limit <= 0:
		return nil, ParamsError.New("limit must be positive")
	case page <= 0:
		return nil, ParamsError.New("page must be positive")
	}

	list, err := service.fundraises.List(ctx, ListParams{
		OrganizerID: creatorID,
		Limit:       limit,
		Page:        page,
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return list, nil
}

// Filled returns funded amount on fundraise by id.
func (service *Service) Filled(ctx context.Context, id uuid.UUID) (float64, error) {
	// TODO: Implement.
	return 0, nil
}
