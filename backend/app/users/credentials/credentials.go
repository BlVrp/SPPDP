package credentials

import (
	"github.com/google/uuid"
)

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
