package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"one-help/app/console/controllers/common"
	"one-help/app/users"
	"one-help/app/users/credentials"
	"one-help/internal/logger"
)

var (
	// ErrUsers is an internal error type for users controller.
	ErrUsers = errs.Class("users controller")
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
		ImageUrl:       request.ImageUrl,
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
		controller.log.Error("error while generating JWT token", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	resp := &AuthResponse{
		User:  ToUserView(user, &credentials.Credentials{Email: request.Email, PhoneNumber: request.PhoneNumber}),
		Token: jwtToken,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
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
// @Failure	400,403,404,500	{object}	common.ErrResponseCode
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
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
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

	creds, err := credentials.GetFromContext(ctx)
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

// GetByID is an endpoint for getting user public info by id.
// @Summary	Provides user public info by id
// @Tags	Users
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Success	200		{object}	UserPublicView
// @Failure 400,404,500	{object}	common.ErrResponseCode
// @Router	/users/{id}	[get].
func (controller *Users) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse id")).Serve(controller.log, ErrUsers, w)
		return
	}

	user, err := controller.users.Get(ctx, userID)
	if err != nil {
		controller.log.Error("failed to get user by id", ErrUsers.Wrap(err))
		if errors.Is(err, users.ErrNoUser) {
			common.NewErrResponse(http.StatusNotFound, users.ErrNoUser).Serve(controller.log, ErrUsers, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get user by id")).Serve(controller.log, ErrUsers, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToUserPublicView(user)); err != nil {
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrUsers, w)
		return
	}
}

// ChangePassword is an endpoint for changing user's password.
// @Summary	Update user's password
// @Tags	Users
// @Produce	json
// @Accept	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Param	request	body	UpdatePasswordRequest	true	"Update password fields"
// @Success	200
// @Failure	401,500	{object}	common.ErrResponseCode
// @Router	/users/change-password	[patch].
func (controller *Users) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creds, err := credentials.GetFromContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		return
	}

	var request UpdatePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode change password request", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrUsers, w)
		return
	}

	err = controller.users.UpdatePassword(ctx, creds.UserID, request.OldPass, request.NewPass)
	if err != nil {
		controller.log.Error("failed to update user's password", ErrUsers.Wrap(err))
		switch {
		case errors.Is(err, users.ErrNoUser):
			common.NewErrResponse(http.StatusNotFound, users.ErrNoUser).Serve(controller.log, ErrUsers, w)
		case users.ParamsError.Has(err):
			common.NewErrResponse(http.StatusForbidden, errors.Unwrap(err)).Serve(controller.log, ErrUsers, w)
		default:
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get user by id")).Serve(controller.log, ErrUsers, w)
		}
		return
	}

	if err = json.NewEncoder(w).Encode(nil); err != nil {
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrUsers, w)
		return
	}
}

// GetRaffleParticipants is an endpoint for getting raffle participants by raffle id.
// @Summary	Provides raffle participants by raffle id
// @Tags	Users
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Success	200		{object}	[]UserPublicView
// @Failure 400,401,404,500	{object}	common.ErrResponseCode
// @Router	/users/raffle-participants/{id}	[get].
func (controller *Users) GetRaffleParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raffleID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse raffle id")).Serve(controller.log, ErrUsers, w)
		return
	}

	userList, err := controller.users.ListRaffleParticipants(ctx, raffleID)
	if err != nil {
		controller.log.Error("failed to get raffle participants", ErrUsers.Wrap(err))
		if errors.Is(err, users.ErrNoUser) {
			common.NewErrResponse(http.StatusNotFound, users.ErrNoUser).Serve(controller.log, ErrUsers, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get raffle participants")).Serve(controller.log, ErrUsers, w)
		return
	}

	var viewList = make([]UserPublicView, len(userList))

	for i, user := range userList {
		viewList[i] = *ToUserPublicView(&user)
	}

	if err = json.NewEncoder(w).Encode(viewList); err != nil {
		controller.log.Error("error while encoding response", ErrUsers.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrUsers, w)
		return
	}
}
