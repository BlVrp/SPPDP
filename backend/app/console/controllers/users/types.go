package users

import (
	"github.com/google/uuid"

	"one-help/app/users"
	"one-help/app/users/credentials"
)

// RegisterRequest defines request values for register endpoint.
type RegisterRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Website        string `json:"website"`
	City           string `json:"city"`
	Post           string `json:"post"`
	PostDepartment string `json:"postDepartment"`
	PhoneNumber    string `json:"phoneNumber"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

// AuthResponse contains user and auth token.
type AuthResponse struct {
	User  *UserView `json:"user"`
	Token string    `json:"token"`
}

// UserView defines user view type.
type UserView struct {
	ID             uuid.UUID `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Website        string    `json:"website"`
	FileName       string    `json:"fileName"`
	City           string    `json:"city"`
	Post           string    `json:"post"`
	PostDepartment string    `json:"postDepartment"`
	PhoneNumber    string    `json:"phoneNumber"`
	Email          string    `json:"email"`
}

// ToUserView builds user view.
func ToUserView(user *users.User, userCreds *credentials.Credentials) *UserView {
	return &UserView{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Website:        user.Website,
		FileName:       user.FileName,
		City:           user.City,
		Post:           user.Post,
		PostDepartment: user.PostDepartment,
		PhoneNumber:    userCreds.PhoneNumber,
		Email:          userCreds.Email,
	}
}

// LoginRequest defines request values for login endpoint.
type LoginRequest struct {
	Identifier string `json:"identifier"` // INFO: Email or phone number.
	Password   string `json:"password"`
}
