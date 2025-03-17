package users

import (
	"github.com/google/uuid"

	"one-help/app/users/auth"
)

// Config defines configuration for users.
type Config struct {
	TokenAuthSecret string `env:"TOKEN_AUTH_SECRET"`
}

// User describes user entity.
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Website   string
	FileName  string

	auth.Credentials

	DeliveryAddresses []DeliveryAddress
}

// FullName returns the full name of the user.
// NOTE: Empty fields will be ignored, empty first and last names will return 'unknown' string.
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return "unknown"
	}

	return u.FirstName + " " + u.LastName
}
