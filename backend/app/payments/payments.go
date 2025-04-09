package payments

import (
	"github.com/google/uuid"
)

// TypeStripe defines stripe payment type.
const TypeStripe string = "STRIPE"

// Payment holds payment info.
type Payment struct {
	DonationId    uuid.UUID
	PaymentType   string
	TransactionId string
	Confirmed     bool
}
