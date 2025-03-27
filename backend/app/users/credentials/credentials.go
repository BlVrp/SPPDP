package credentials

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ctxKeyType string

// credsKey defines key for context store.
const credsKey ctxKeyType = "credentials"

// Credentials holds user's credential data.
type Credentials struct {
	UserID       uuid.UUID `json:"userId,omitempty"`
	PhoneNumber  string    `json:"phoneNumber,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"-"`
}

// IsEmpty returns true if phone number and email are empty.
func (c *Credentials) IsEmpty() bool {
	return c.Email == "" && c.PhoneNumber == ""
}

// SetIntoContext returns updates context with credentials set inside.
func (c *Credentials) SetIntoContext(ctx context.Context) context.Context {
	return SetIntoContext(ctx, c)
}

// SetIntoContext returns updates context with provided credentials set inside.
func SetIntoContext(ctx context.Context, creds *Credentials) context.Context {
	return context.WithValue(ctx, credsKey, creds)
}

// GetIntoContext returns credentials from provided context if exists.
func GetIntoContext(ctx context.Context) (*Credentials, error) {
	val := ctx.Value(credsKey)
	if val == nil {
		return nil, errors.New("credentials not found in context")
	}

	creds, ok := val.(*Credentials)
	if !ok {
		return nil, errors.New("value is not assignable to credentials type")
	}

	return creds, nil
}
