package fundraises

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"one-help/app/donations"
	"one-help/app/fundraises/statuses"
	"one-help/app/payments"
	"one-help/app/stripe"
	"one-help/internal/logger"
)

var (
	// Error wraps errors from fundraises service that indicates about internal errors.
	Error = errs.Class("fundraises service")
	// ParamsError wraps errors from fundraises service that indicates about invalid or malformed parameters' data.
	ParamsError = errs.Class("fundraises service: params")
)

// Service handles fundraises related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger

	fundraises DB
	donations  donations.DB
	payments   payments.DB

	charger *stripe.Charger
}

// NewService is a constructor for fundraises service.
func NewService(logger logger.Logger, fundraises DB, donations donations.DB, payments payments.DB, charger *stripe.Charger) *Service {
	return &Service{
		logger:     logger,
		fundraises: fundraises,
		donations:  donations,
		payments:   payments,
		charger:    charger,
	}
}

// Create creates new Fundraise data in the system.
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
		ImageUrl:     params.ImageUrl,
	}

	err := service.fundraises.Create(ctx, *fundraise)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return fundraise, nil
}

// Get returns Fundraise by id.
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
	filled, err := service.fundraises.GetFilled(ctx, id)
	if err != nil {
		return 0, Error.Wrap(err)
	}

	return filled, nil
}

// RegisterDonate register new donate values, provides payment url.
func (service *Service) RegisterDonate(ctx context.Context, params RegisterDonateParams) (result RegisterDonateResult, err error) {
	donation := donations.Donation{
		ID:          uuid.New(),
		UserId:      params.UserID,
		FundraiseId: params.FundraiseID,
		Amount:      0,
		CreatedAt:   time.Now().UTC(),
	}
	err = service.donations.Create(ctx, donation)
	if err != nil {
		return result, Error.Wrap(err)
	}

	var paymentID string
	result.PaymentURL, paymentID, err = service.charger.SetupChargeSessionURL("/fundraises/donations/" + donation.ID.String())
	if err != nil {
		return result, Error.Wrap(err)
	}

	err = service.payments.Create(ctx, payments.Payment{
		DonationId:    donation.ID,
		PaymentType:   payments.TypeStripe,
		TransactionId: paymentID,
		Confirmed:     false,
	})
	if err != nil {
		return result, Error.Wrap(err)
	}

	return result, nil
}

// GetDonation returns donation by id.
func (service *Service) GetDonation(ctx context.Context, donationID uuid.UUID) (donation donations.Donation, err error) {
	donation, err = service.donations.Get(ctx, donationID)
	if err != nil {
		return donation, Error.Wrap(err)
	}

	return donation, nil
}

// ConfirmDonation finishes donation processes.
func (service *Service) ConfirmDonation(ctx context.Context, donation donations.Donation) error {
	payment, err := service.payments.Get(ctx, donation.ID)
	if err != nil {
		return Error.Wrap(err)
	}

	paid, funded, err := service.charger.GetSessionPaymentData(payment.TransactionId)
	if err != nil {
		return Error.Wrap(err)
	}

	if !paid {
		service.logger.WarnF("receiver unpaid session: %s, for donation: %s in ConfirmDonation", payment.TransactionId, donation.ID.String())
	}

	donation.Amount = funded
	payment.Confirmed = true

	err = service.donations.Update(ctx, donation)
	if err != nil {
		return Error.Wrap(err)
	}

	err = service.payments.Update(ctx, payment)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// CancelDonation removes canceled donation data.
func (service *Service) CancelDonation(ctx context.Context, donation donations.Donation) error {
	return Error.Wrap(service.donations.Delete(ctx, donation.ID))
}
