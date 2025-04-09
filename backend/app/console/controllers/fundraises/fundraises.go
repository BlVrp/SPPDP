package fundraises

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"one-help/app/console/controllers/common"
	"one-help/app/fundraises"
	"one-help/app/users/credentials"
	"one-help/internal/logger"
)

var (
	// ErrFundraises is an internal error type for fundraises controller.
	ErrFundraises = errs.Class("fundraises controller")
)

// Fundraises is a controller that handles all fundraises related routes.
type Fundraises struct {
	log logger.Logger

	fundraises *fundraises.Service

	frontEndRedirectUrl string
}

// NewFundraises is a constructor for fundraises controller.
func NewFundraises(log logger.Logger, fundraises *fundraises.Service, frontEndRedirectUrl string) *Fundraises {
	fundraisesController := &Fundraises{
		log:                 log,
		fundraises:          fundraises,
		frontEndRedirectUrl: frontEndRedirectUrl,
	}

	return fundraisesController
}

// Create is an endpoint for creating new fundraise.
// @Summary	Creates new fundraise
// @Tags	Fundraises
// @Accept	json
// @Produce	json
// @Param	request	body	CreateRequest	true	"Fundraise needed data fields"
// @Success	200		{object}	FundraiseView
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/fundraises/	[post].
func (controller *Fundraises) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode create request body", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrFundraises, w)
		return
	}

	// INFO: Caller creds.
	creds, err := credentials.GetFromContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, err).Serve(controller.log, ErrFundraises, w)
		return
	}

	createParams := fundraises.CreateParams{
		OrganizerId:  creds.UserID,
		Title:        request.Title,
		Description:  request.Description,
		TargetAmount: request.TargetAmount,
		EndDate:      request.EndDate,
		ImageUrl:     request.ImageUrl,
	}

	fundraise, err := controller.fundraises.Create(ctx, createParams)
	if err != nil {
		controller.log.Error("error while creating fundraise", ErrFundraises.Wrap(err))
		if fundraises.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrFundraises, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrFundraises, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToFundraiseView(fundraise, 0)); err != nil {
		controller.log.Error("error while encoding response:", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrFundraises, w)
		return
	}
}

// List is an endpoint for listing fundraises.
// @Summary	Returns list of fundraises
// @Tags	Fundraises
// @Produce	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Param	limit			query	integer	false	"Items per page (positive number expected) [default value: 20]"
// @Param	page			query	integer	false	"Number of the page (1...) [default value: 1]"
// @Success	200		{object}	common.Page[FundraiseView]
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/fundraises/	[get].
func (controller *Fundraises) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		limit = 20
		page  = 1
		err   error
	)
	if val := r.URL.Query().Get("limit"); val != "" {
		limit, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'limit' query parameter", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid limit value")).Serve(controller.log, ErrFundraises, w)
			return
		}
	}
	if val := r.URL.Query().Get("page"); val != "" {
		page, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'page' query parameter", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid page value")).Serve(controller.log, ErrFundraises, w)
			return
		}
	}

	list, err := controller.fundraises.List(ctx, limit, page, nil)
	if err != nil {
		controller.log.Error("failed to list fundraises", ErrFundraises.Wrap(err))
		if fundraises.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrFundraises, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list fundraises")).Serve(controller.log, ErrFundraises, w)
		return
	}

	var filled float64
	var viewList = make([]*FundraiseView, len(list))
	for i, fundraise := range list {
		filled, err = controller.fundraises.Filled(ctx, fundraise.ID)
		if err != nil {
			controller.log.Error("failed to list fundraises", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list fundraises")).Serve(controller.log, ErrFundraises, w)
			return
		}

		viewList[i] = ToFundraiseView(&fundraise, filled)
	}

	resp := &common.Page[*FundraiseView]{
		Data:  viewList,
		Page:  page,
		Limit: limit,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrFundraises, w)
		return
	}
}

// GetByID is an endpoint for getting fundraise by id.
// @Summary	Provides fundraise by id
// @Tags	Fundraises
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Success	200		{object}	FundraiseView
// @Failure 400,401,404,500	{object}	common.ErrResponseCode
// @Router	/fundraises/{id}	[get].
func (controller *Fundraises) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fundsraiseID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse id")).Serve(controller.log, ErrFundraises, w)
		return
	}

	fundraise, err := controller.fundraises.Get(ctx, fundsraiseID)
	if err != nil {
		controller.log.Error("failed to get fundraise by id", ErrFundraises.Wrap(err))
		if errors.Is(err, fundraises.ErrNoFundraise) {
			common.NewErrResponse(http.StatusNotFound, fundraises.ErrNoFundraise).Serve(controller.log, ErrFundraises, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get fundraise by id")).Serve(controller.log, ErrFundraises, w)
		return
	}

	filled, err := controller.fundraises.Filled(ctx, fundraise.ID)
	if err != nil {
		controller.log.Error("failed to get fundraise filled value", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get fundraise filled value")).Serve(controller.log, ErrFundraises, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToFundraiseView(fundraise, filled)); err != nil {
		controller.log.Error("error while encoding response", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrFundraises, w)
		return
	}
}

// ListMy is an endpoint for listing user's fundraises.
// @Summary	Returns list of fundraises of the user by bearer token.
// @Tags	Fundraises
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Param	limit			query	integer	false	"Items per page (positive number expected) [default value: 20]"
// @Param	page			query	integer	false	"Number of the page (1...) [default value: 1]"
// @Success	200		{object}	common.Page[FundraiseView]
// @Failure 400,401,404,500	{object}	common.ErrResponseCode
// @Router	/fundraises/my	[get].
func (controller *Fundraises) ListMy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// INFO: Caller creds.
	creds, err := credentials.GetFromContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, err).Serve(controller.log, ErrFundraises, w)
		return
	}

	var (
		limit = 20
		page  = 1
	)
	if val := r.URL.Query().Get("limit"); val != "" {
		limit, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'limit' query parameter", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid limit value")).Serve(controller.log, ErrFundraises, w)
			return
		}
	}
	if val := r.URL.Query().Get("page"); val != "" {
		page, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'page' query parameter", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid page value")).Serve(controller.log, ErrFundraises, w)
			return
		}
	}

	list, err := controller.fundraises.List(ctx, limit, page, &creds.UserID)
	if err != nil {
		controller.log.Error("failed to list creator fundraises", ErrFundraises.Wrap(err))
		if fundraises.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrFundraises, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list creator fundraises")).Serve(controller.log, ErrFundraises, w)
		return
	}

	var filled float64
	var viewList = make([]*FundraiseView, len(list))
	for i, fundraise := range list {
		filled, err = controller.fundraises.Filled(ctx, fundraise.ID)
		if err != nil {
			controller.log.Error("failed to list creator fundraises", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list creator fundraises")).Serve(controller.log, ErrFundraises, w)
			return
		}

		viewList[i] = ToFundraiseView(&fundraise, filled)
	}

	resp := &common.Page[*FundraiseView]{
		Data:  viewList,
		Page:  page,
		Limit: limit,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrFundraises, w)
		return
	}
}

// Donate is an endpoint creating new payment url to donate to fundraise.
// @Summary	Creates new payment url to donate to fundraise.
// @Tags	Fundraises
// @Produce	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Success	200	{object}	DonateResponse
// @Failure	400,401,404,500	{object}	common.ErrResponseCode
// @Router	/fundraises/{id}/donate	[post].
func (controller *Fundraises) Donate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// INFO: Caller creds.
	creds, err := credentials.GetFromContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, err).Serve(controller.log, ErrFundraises, w)
		return
	}

	fundsraiseID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse id")).Serve(controller.log, ErrFundraises, w)
		return
	}

	fundraise, err := controller.fundraises.Get(ctx, fundsraiseID)
	if err != nil {
		controller.log.Error("failed to get fundraise by id", ErrFundraises.Wrap(err))
		if errors.Is(err, fundraises.ErrNoFundraise) {
			common.NewErrResponse(http.StatusNotFound, fundraises.ErrNoFundraise).Serve(controller.log, ErrFundraises, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get fundraise by id")).Serve(controller.log, ErrFundraises, w)
		return
	}

	result, err := controller.fundraises.RegisterDonate(ctx, fundraises.RegisterDonateParams{
		FundraiseID: fundraise.ID,
		UserID:      creds.UserID,
	})
	if err != nil {
		controller.log.Error("failed to register payment", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to register payment")).Serve(controller.log, ErrFundraises, w)
		return
	}

	if err = json.NewEncoder(w).Encode(&DonateResponse{PaymentURL: result.PaymentURL}); err != nil {
		controller.log.Error("error while encoding response", ErrFundraises.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrFundraises, w)
		return
	}
}

// FinishDonation is an endpoint payment callback.
// @Summary	Finishes payment process (for internal use).
// @Tags	Donations
// @Produce	json
// @Success	301
// @Failure	400,404,500	{object}	common.ErrResponseCode
// @Router	/fundraises/donations/{id}	[get].
func (controller *Fundraises) FinishDonation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	donationID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse donation id")).Serve(controller.log, ErrFundraises, w)
		return
	}

	switch {
	case r.URL.Query().Get("success") == "true":
		if err = controller.fundraises.ConfirmDonation(ctx, donationID); err != nil {
			controller.log.Error("failed to confirm donation", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to confirm donation")).Serve(controller.log, ErrFundraises, w)
			return
		}

		http.Redirect(w, r, controller.frontEndRedirectUrl+"?success=true", http.StatusMovedPermanently)
	case r.URL.Query().Get("canceled") == "true":
		if err = controller.fundraises.CancelDonation(ctx, donationID); err != nil {
			controller.log.Error("failed to cancel donation", ErrFundraises.Wrap(err))
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to cancel donation")).Serve(controller.log, ErrFundraises, w)
			return
		}

		http.Redirect(w, r, controller.frontEndRedirectUrl+"?canceled=true", http.StatusMovedPermanently)
	default:
		controller.log.Warn("invalid callback url format received")
		common.NewErrResponse(http.StatusBadRequest, errors.New("failed to cancel donation")).Serve(controller.log, ErrFundraises, w)
	}
}
