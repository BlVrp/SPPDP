package users

import (
	"github.com/google/uuid"
)

// Config defines configuration for users.
type Config struct {
	TokenAuthSecret   string `env:"TOKEN_AUTH_SECRET"`
	EmailRegExp       string `env:"EMAIL_REGEXP"`
	PhoneNumberRegExp string `env:"PHONE_NUMBER_REGEXP"`
}

// User describes user entity.
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Website   string
	ImageUrl  string

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

// RegisterParams holds parameters needed to register new user.
type RegisterParams struct {
	FirstName string
	LastName  string
	Website   string
	ImageUrl  string

	City           string
	Post           string
	PostDepartment string

	PhoneNumber string
	Email       string
	Password    string
}

type UserWithContacts struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Website   string
	ImageUrl  string

	City           string
	Post           string
	PostDepartment string

	PhoneNumber string
	Email       string
}

// AuthorizeParams holds parameters needed to authorize user.
type AuthorizeParams struct {
	Identifier string // INFO: email of phone number.
	Password   string
}
