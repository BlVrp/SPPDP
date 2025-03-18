package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoUserCredentials indicates that user's credentials does not exist.
	ErrNoUserCredentials = errs.New("user credentials does not exist")
	// ErrUserAlreadyHasCredentials indicates that user's credentials already exists.
	ErrUserAlreadyHasCredentials = errs.New("user credentials already exists")
	// ErrUserEmailTaken indicates that such email already exists.
	ErrUserEmailTaken = errs.New("email is already taken")
	// ErrUserPhoneNumberTaken indicates that such phone number already exists.
	ErrUserPhoneNumberTaken = errs.New("phone number is already taken")
)

// DB exposes access to user credentials db.
//
// architecture: DB
type DB interface {
	// Create inserts user's credentials into the database.
	Create(ctx context.Context, creds Credentials) error
	// Get returns user's credentials from the database by user's ID.
	Get(ctx context.Context, key GetKey) (Credentials, error)
	// Update updates user's credentials in database by id.
	Update(ctx context.Context, creds Credentials) error
}

// GetKey defines key selector for get query.
type GetKey struct {
	fieldName string
	value     any
}

// FieldName returns the name of the key field.
func (gk *GetKey) FieldName() string {
	return gk.fieldName
}

// Value returns the value of the key.
func (gk *GetKey) Value() any {
	return gk.value
}

// NewGetByID builds get query selector by ID.
func NewGetByID(id uuid.UUID) GetKey {
	return GetKey{
		fieldName: "user_id",
		value:     id,
	}
}

// NewGetByPhoneNumber builds get query selector by phone number.
func NewGetByPhoneNumber(phoneNumber string) GetKey {
	return GetKey{
		fieldName: "phone_number",
		value:     phoneNumber,
	}
}

// NewGetByEmail builds get query selector by email.
func NewGetByEmail(email string) GetKey {
	return GetKey{
		fieldName: "email",
		value:     email,
	}
}
