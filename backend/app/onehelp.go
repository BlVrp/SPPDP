package app

import (
	"one-help/app/donations"
	"one-help/app/events"
	eventformats "one-help/app/events/formats"
	eventparticipants "one-help/app/events/participants"
	eventstatuses "one-help/app/events/statuses"
	"one-help/app/fundraises"
	fundraisestatuses "one-help/app/fundraises/statuses"
	"one-help/app/payments"
	paymenttypes "one-help/app/payments/types"
	"one-help/app/posts"
	"one-help/app/raffles"
	"one-help/app/users"
	"one-help/app/users/credentials"
)

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// Users provides access to users.DB.
	Users() users.DB

	// Credentials provides access to credentials.DB.
	Credentials() credentials.DB

	// Posts provides access to posts.DB.
	Posts() posts.DB

	// FundraiseStatuses provides access to fundraise statuses DB.
	FundraiseStatuses() fundraisestatuses.DB

	// Fundraises provides access to fundraises DB.
	Fundraises() fundraises.DB

	// Events provides access to events DB.
	Events() events.DB

	// EventStatuses provides access to event statuses DB.
	EventStatuses() eventstatuses.DB

	// EventFormats provides access to event formats DB.
	EventFormats() eventformats.DB

	// EventParticipants provides access to event participants DB.
	EventParticipants() eventparticipants.DB

	// Donations provides access to donations DB.
	Donations() donations.DB

	// PaymentTypes provides access to payment types DB.
	PaymentTypes() paymenttypes.DB

	// Payments provides access to payments DB.
	Payments() payments.DB

	// Raffles provides access to raffles DB.
	Raffles() raffles.DB

	// ExecuteMigrations applies migrations for the database.
	ExecuteMigrations(migrationsPath string, isUp bool) (err error)
}
