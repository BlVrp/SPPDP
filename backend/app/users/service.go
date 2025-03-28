package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"one-help/app/users/credentials"
	"one-help/internal/jwt"
	"one-help/internal/logger"
)

var (
	// Error wraps errors from users service that indicates about internal errors.
	Error = errs.Class("users service")
	// ParamsError wraps errors from users service that indicates about invalid or malformed parameters data.
	ParamsError = errs.Class("users service: params")
)

// Service handles users related logic.
//
// architecture: Service
type Service struct {
	logger logger.Logger
	config Config

	tokenizer jwt.Tokenizer[credentials.Credentials]

	users       DB
	credentials credentials.DB

	emailChecker *regexp.Regexp
	phoneChecker *regexp.Regexp
}

// NewService is a constructor for users service.
func NewService(logger logger.Logger, config Config, users DB, creds credentials.DB) *Service {
	return &Service{
		logger:       logger,
		config:       config,
		users:        users,
		credentials:  creds,
		tokenizer:    jwt.NewHS256[credentials.Credentials]([]byte(config.TokenAuthSecret)),
		emailChecker: regexp.MustCompile(config.EmailRegExp),
		phoneChecker: regexp.MustCompile(config.PhoneNumberRegExp),
	}
}

// Register created new User data in the system.
func (service *Service) Register(ctx context.Context, params RegisterParams) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Website:   params.Website,
		FileName:  params.FileName,
		DeliveryAddress: DeliveryAddress{
			City:           params.City,
			Post:           params.Post,
			PostDepartment: params.PostDepartment,
		},
	}
	if err := service.verifyUserData(user); err != nil {
		return nil, err
	}

	creds := &credentials.Credentials{
		UserID:       user.ID,
		PhoneNumber:  params.PhoneNumber,
		Email:        params.Email,
		PasswordHash: service.hashPassword(params.Password),
	}
	if err := service.verifyCredentialData(creds); err != nil {
		return nil, err
	}

	if err := service.users.Create(ctx, *user); err != nil {
		return nil, Error.Wrap(err)
	}

	if err := service.credentials.Create(ctx, *creds); err != nil {
		switch {
		case errors.Is(err, credentials.ErrUserPhoneNumberTaken):
			return nil, ParamsError.Wrap(credentials.ErrUserPhoneNumberTaken)
		case errors.Is(err, credentials.ErrUserEmailTaken):
			return nil, ParamsError.Wrap(credentials.ErrUserEmailTaken)
		}

		return nil, Error.Wrap(err)
	}

	return user, nil
}

// Authorize authorizes user by credentials.
func (service *Service) Authorize(ctx context.Context, params AuthorizeParams) (*User, error) {
	var key credentials.GetKey
	switch {
	case service.phoneChecker.MatchString(params.Identifier):
		key = credentials.NewGetByPhoneNumber(params.Identifier)
	case service.emailChecker.MatchString(params.Identifier):
		key = credentials.NewGetByEmail(params.Identifier)
	default:
		return nil, ParamsError.New("invalid identifier is provided")
	}

	creds, err := service.credentials.Get(ctx, key)
	if err != nil {
		if errors.Is(err, credentials.ErrNoUserCredentials) {
			return nil, ParamsError.Wrap(ErrNoUser)
		}

		return nil, Error.Wrap(err)
	}

	if creds.PasswordHash != service.hashPassword(params.Password) {
		return nil, ParamsError.Wrap(ErrInvalidPassword)
	}

	user, err := service.users.Get(ctx, creds.UserID)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &user, nil
}

// Update updates all User data.
func (service *Service) Update(ctx context.Context, user User) error {
	_, err := service.users.Get(ctx, user.ID)
	if err != nil {
		return Error.Wrap(err)
	}

	if err = service.verifyUserData(&user); err != nil {
		return err
	}

	if err = service.users.Update(ctx, user); err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// Get returns User by ID.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := service.users.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoUser) {
			return nil, ParamsError.Wrap(ErrNoUser)
		}

		return nil, Error.Wrap(err)
	}

	return &user, nil
}

// GetCreds returns User creds by ID.
func (service *Service) GetCreds(ctx context.Context, id uuid.UUID) (*credentials.Credentials, error) {
	creds, err := service.credentials.Get(ctx, credentials.NewGetByID(id))
	if err != nil {
		if errors.Is(err, credentials.ErrNoUserCredentials) {
			return nil, ParamsError.Wrap(ErrNoUser)
		}

		return nil, Error.Wrap(err)
	}

	return &creds, nil
}

// UpdatePassword updated user password.
func (service *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error {
	creds, err := service.credentials.Get(ctx, credentials.NewGetByID(userID))
	if err != nil {
		if errors.Is(err, credentials.ErrNoUserCredentials) {
			return ParamsError.Wrap(ErrNoUser)
		}

		return Error.Wrap(err)
	}

	if creds.PasswordHash != service.hashPassword(oldPass) {
		return ParamsError.New("invalid old password")
	}

	creds.PasswordHash = service.hashPassword(newPass)
	err = service.credentials.Update(ctx, creds)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// JWTToken generates JWT token for user.
func (service *Service) JWTToken(ctx context.Context, userID uuid.UUID) (string, error) {
	creds, err := service.credentials.Get(ctx, credentials.NewGetByID(userID))
	if err != nil {
		if errors.Is(err, credentials.ErrNoUserCredentials) {
			return "", ParamsError.Wrap(ErrNoUser)
		}

		return "", Error.Wrap(err)
	}

	token, err := service.tokenizer.Token(creds)
	if err != nil {
		return "", Error.Wrap(err)
	}

	tokenStr, err := token.String()
	if err != nil {
		return "", Error.Wrap(err)
	}

	return tokenStr, nil
}

// ValidateJWTToken validates provided token and user existence.
func (service *Service) ValidateJWTToken(ctx context.Context, tkn string) (*jwt.Token[credentials.Credentials], error) {
	token := new(jwt.Token[credentials.Credentials])
	if err := token.Parse(tkn); err != nil {
		return nil, Error.Wrap(err)
	}

	if err := service.tokenizer.Verify(token); err != nil {
		return nil, Error.Wrap(err)
	}

	_, err := service.credentials.Get(ctx, credentials.NewGetByID(token.Payload.UserID))
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return token, nil
}

// SetCredsToContext returns updated context with user credentials fetched by user id.
func (service *Service) SetCredsToContext(ctx context.Context, id uuid.UUID) (context.Context, error) {
	creds, err := service.GetCreds(ctx, id)
	if err != nil {
		return ctx, err
	}

	return credentials.SetIntoContext(ctx, creds), nil
}

// verifyUserData returns error in case of incorrect user data.
func (service *Service) verifyUserData(user *User) error {
	switch {
	case user.FirstName == "":
		return ParamsError.New("first name is empty")
	case user.LastName == "":
		return ParamsError.New("last name is empty")
	case user.IsDeliveryAddressEmpty():
		return ParamsError.New("delivery address is required")
	}

	return nil
}

// verifyCredentialData returns error in case of incorrect user credential data.
func (service *Service) verifyCredentialData(creds *credentials.Credentials) error {
	switch {
	case creds.IsEmpty():
		return ParamsError.New("email or phone number is required")
	case !service.phoneChecker.MatchString(creds.PhoneNumber):
		return ParamsError.New("invalid phone number format %s |", creds.PhoneNumber)
	case !service.emailChecker.MatchString(creds.Email):
		return ParamsError.New("invalid email address format")
	}

	return nil
}

// hashPassword returns hash of the provided password in hex.
func (service *Service) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))

	return hex.EncodeToString(hash[:])
}
