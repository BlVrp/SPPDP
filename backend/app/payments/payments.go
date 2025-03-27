package payments

import "github.com/google/uuid"

type Payment struct {
	DonationId    uuid.UUID
	PaymentType   string
	TransactionId uuid.UUID
	Confirmed     bool
}
