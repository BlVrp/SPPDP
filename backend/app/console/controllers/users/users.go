package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zeebo/errs"

	"one-help/app/console/controllers/common"
	"one-help/app/users"
	"one-help/app/users/credentials"
	"one-help/internal/logger"
)

var (
	// ErrUsers is an internal error type for users controller.
	ErrUsers = errs.Class("users controller error")
)

// Users is a controller that handles all users related routes.
type Users struct {
	log logger.Logger

	users *users.Service
}

// NewUsers is a constructor for users controller.
func NewUsers(log logger.Logger, users *users.Service) *Users {
	usersController := &Users{
		log:   log,
		users: users,
	}

	return usersController
}

// Register is an endpoint for creating new user.
// @Summary	Register new user
// @Tags	Auth
// @Accept	json
// @Produce	json
// @Param	request	body	RegisterRequest	true	"Register request fields"
// @Success	200		{object}	AuthResponse
// @Failure	400,500	{object}	common.ErrResponseCode
// @Router	/auth/register	[post].
func (controller *Users) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode register request", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrUsers, w)
		return
	}

	registerParams := users.RegisterParams{
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Website:        request.Website,
		City:           request.City,
		Post:           request.Post,
		PostDepartment: request.PostDepartment,
		PhoneNumber:    request.PhoneNumber,
		Email:          request.Email,
		Password:       request.Password,
	}

	user, err := controller.users.Register(ctx, registerParams)
	if err != nil {
		controller.log.Error("error while register:", ErrUsers.Wrap(err))
		if users.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	jwtToken, err := controller.users.JWTToken(ctx, user.ID)
	if err != nil {
		controller.log.Error("error while generating JWT token:", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	resp := &AuthResponse{
		User:  ToUserView(user, &credentials.Credentials{Email: request.Email, PhoneNumber: request.PhoneNumber}),
		Token: jwtToken,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response:", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}
}

// Login is an endpoint for logging into the system.
// @Summary	Login user
// @Tags	Auth
// @Accept	json
// @Produce	json
// @Param	request	body	LoginRequest	true	"Login request fields"
// @Success	200			{object}	AuthResponse
// @Failure	400,404,500	{object}	common.ErrResponseCode
// @Router	/auth/login	[post].
func (controller *Users) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode login request", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrUsers, w)
		return
	}

	user, err := controller.users.Authorize(ctx, users.AuthorizeParams{
		Identifier: request.Identifier,
		Password:   request.Password,
	})
	if err != nil {
		controller.log.Error("filed to authorize user", ErrUsers.Wrap(err))
		switch {
		case errors.Is(err, users.ErrNoUser):
			common.NewErrResponse(http.StatusNotFound, users.ErrNoUser).Serve(controller.log, ErrUsers, w)
		case errors.Is(err, users.ErrInvalidPassword):
			common.NewErrResponse(http.StatusForbidden, users.ErrInvalidPassword).Serve(controller.log, ErrUsers, w)
		case users.ParamsError.Has(err):
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		default:
			common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		}
		return
	}

	jwtToken, err := controller.users.JWTToken(ctx, user.ID)
	if err != nil {
		controller.log.Error("error while generating JWT token", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	creds, err := controller.users.GetCreds(ctx, user.ID)
	if err != nil {
		controller.log.Error("error while getting creds", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	resp := &AuthResponse{
		User:  ToUserView(user, creds),
		Token: jwtToken,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response:", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}
}

// Get is an endpoint for getting user info from Authorization token.
// @Summary	Get user from Authorization token
// @Tags	Users
// @Produce	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Success	200		{object}	UserView
// @Failure	401,500	{object}	common.ErrResponseCode
// @Router	/users/	[get].
func (controller *Users) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creds, err := credentials.GetIntoContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	user, err := controller.users.Get(ctx, creds.UserID)
	if err != nil {
		controller.log.Error("failed to get user by id", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get user by id")).Serve(controller.log, ErrUsers, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToUserView(user, creds)); err != nil {
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrUsers, w)
		return
	}
}
