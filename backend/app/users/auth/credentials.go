package auth

import (
	"github.com/google/uuid"
)

// Credentials holds user's credential data.
type Credentials struct {
	UserID       uuid.UUID
	PhoneNumber  string
	Email        string
	PasswordHash string
}
