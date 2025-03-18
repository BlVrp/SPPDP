package users

import (
	"github.com/google/uuid"
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

	DeliveryAddress
}

// FullName returns the full name of the user.
// NOTE: Empty fields will be ignored, empty first and last names will return 'unknown' string.
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return "unknown"
	}

	return u.FirstName + " " + u.LastName
}

// IsDeliveryAddressFull returns true if delivery address is full-filled.
func (u *User) IsDeliveryAddressFull() bool {
	return u.City != "" && u.Post != "" && u.PostDepartment != ""
}

// IsDeliveryAddressEmpty returns true if delivery address is full-empty.
func (u *User) IsDeliveryAddressEmpty() bool {
	return u.City == "" && u.Post == "" && u.PostDepartment == ""
}

// EmptyDeliveryAddress clear delivery address.
func (u *User) EmptyDeliveryAddress() {
	u.DeliveryAddress = DeliveryAddress{}
}
