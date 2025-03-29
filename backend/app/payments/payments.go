package payments

import "github.com/google/uuid"

// Payment holds payment info.
type Payment struct {
	DonationId    uuid.UUID
	PaymentType   string
	TransactionId string
	Confirmed     bool
}
