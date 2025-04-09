package stripe

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/zeebo/errs"

	"one-help/internal/logger"
)

// Error is an error wrapper that notifies that error was produced by stripe charger.
var Error = errs.Class("stripe charger")

// Config holds configurable values for Stripe charger.
type Config struct {
	SecretAPIKey   string `env:"SECRET_API_KEY"`
	RedirectDomain string `env:"REDIRECT_DOMAIN"` // INFO: Ex.: localhost:port/api/v0
	PriceID        string `env:"PRICE_ID"`
}

// Charger defines Stripe charger functionality.
type Charger struct {
	log    logger.Logger
	config Config
}

// NewCharger is a constructor for Charger.
func NewCharger(log logger.Logger, config Config) *Charger {
	stripe.Key = config.SecretAPIKey
	return &Charger{
		log:    log,
		config: config,
	}
}

// SetupChargeSessionURL setups charge session and provides payment redirect url.
// NOTE: redirectPath must start with '/'.
func (c *Charger) SetupChargeSessionURL(redirectPath string) (url string, id string, err error) {
	redirectPath = c.config.RedirectDomain + redirectPath
	checkoutParams := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Price:    stripe.String(c.config.PriceID),
			Quantity: stripe.Int64(1),
		}},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(redirectPath + "?success=true"),
		CancelURL:  stripe.String(redirectPath + "?canceled=true"),
		Currency:   stripe.String(string(stripe.CurrencyUAH)),
	}

	session_, err := session.New(checkoutParams)
	if err != nil {
		c.log.Error("error creating session", Error.Wrap(err))
		return url, id, Error.Wrap(err)
	}

	url = session_.URL
	id = session_.ID

	return url, id, nil
}

// GetSessionPaymentData provides payment amount info.
func (c *Charger) GetSessionPaymentData(id string) (paid bool, funded float64, err error) {
	session_, err := session.Get(id, new(stripe.CheckoutSessionParams))
	if err != nil {
		c.log.Error("error getting session", Error.Wrap(err))
		return paid, funded, Error.Wrap(err)
	}

	paid = session_.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid
	funded = float64(session_.AmountTotal) / 100.

	return paid, funded, nil
}
