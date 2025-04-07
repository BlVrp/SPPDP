package raffles

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"one-help/app/console/controllers/common"
	"one-help/app/raffles"
	"one-help/internal/logger"
)

var (
	// ErrRaffles is an internal error type for raffles controller.
	ErrRaffles = errs.Class("raffles controller")
)

// Raffles is a controller that handles all raffles related routes.
type Raffles struct {
	log logger.Logger

	raffles *raffles.Service
}

// NewRaffles is a constructor for raffles controller.
func NewRaffles(log logger.Logger, raffles *raffles.Service) *Raffles {
	rafflesController := &Raffles{
		log:     log,
		raffles: raffles,
	}

	return rafflesController
}

// Create is an endpoint for creating new raffle.
// @Summary	Creates new raffle
// @Tags	Raffles
// @Accept	json
// @Produce	json
// @Param	request	body	CreateRequest	true	"Raffle needed data fields"
// @Success	200		{object}	RaffleView
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/raffles/	[post].
func (controller *Raffles) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode create request body", ErrRaffles.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrRaffles, w)
		return
	}

	var giftsParams = make([]raffles.GiftCreateParams, len(request.Gifts))
	for i, gift := range request.Gifts {
		giftsParams[i] = raffles.GiftCreateParams{
			Title:       gift.Title,
			Description: gift.Description,
			ImageLink:   gift.ImageLink,
		}
	}

	createParams := raffles.CreateParams{
		Title:           request.Title,
		Description:     request.Description,
		MinimumDonation: request.MinimumDonation,
		StartDate:       request.StartDate,
		EndDate:         request.EndDate,
		FundraiseID:     request.FundraiseID,
		Gifts:           giftsParams,
	}

	raffle, gifts, err := controller.raffles.Create(ctx, createParams)
	if err != nil {
		controller.log.Error("error while creating raffle", ErrRaffles.Wrap(err))
		if raffles.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrRaffles, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrRaffles, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToRaffleView(raffle, gifts)); err != nil {
		controller.log.Error("error while encoding response:", ErrRaffles.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrRaffles, w)
		return
	}
}

// List is an endpoint for listing raffles.
// @Summary	Returns list of raffles
// @Tags	Raffles
// @Produce	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Param	limit			query	integer	false	"Items per page (positive number expected) [default value: 20]"
// @Param	page			query	integer	false	"Number of the page (1...) [default value: 1]"
// @Param	fundraiseId		query	uuid	false	"Optional filtration by fundraise"
// @Success	200		{object}	common.Page[RaffleView]
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/raffles/	[get].
func (controller *Raffles) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		fundraiseID *uuid.UUID
		limit       = 20
		page        = 1
		err         error
	)
	if val := r.URL.Query().Get("limit"); val != "" {
		limit, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'limit' query parameter", ErrRaffles.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid limit value")).Serve(controller.log, ErrRaffles, w)
			return
		}
	}
	if val := r.URL.Query().Get("page"); val != "" {
		page, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'page' query parameter", ErrRaffles.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid page value")).Serve(controller.log, ErrRaffles, w)
			return
		}
	}
	if val := r.URL.Query().Get("fundraiseId"); val != "" {
		var fundraiseID_ uuid.UUID
		fundraiseID_, err = uuid.Parse(val)
		if err != nil {
			controller.log.Error("failed to parse 'fundraiseId' query parameter", ErrRaffles.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid fundraiseId value")).Serve(controller.log, ErrRaffles, w)
			return
		}

		fundraiseID = &fundraiseID_
	}

	list, err := controller.raffles.List(ctx, limit, page, fundraiseID)
	if err != nil {
		controller.log.Error("failed to list raffles", ErrRaffles.Wrap(err))
		if raffles.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrRaffles, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list raffles")).Serve(controller.log, ErrRaffles, w)
		return
	}

	resp := &common.Page[*RaffleView]{
		Data:  make([]*RaffleView, len(list)),
		Page:  page,
		Limit: limit,
	}

	for i, raffle := range list {
		var gifts []raffles.Gift
		gifts, err = controller.raffles.ListGifts(ctx, raffle.FundraiseID)
		if err != nil {
			controller.log.Error("error while listing gifts", ErrRaffles.Wrap(err))
			common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list raffle gifts")).Serve(controller.log, ErrRaffles, w)
			return
		}
		resp.Data[i] = ToRaffleView(&raffle, gifts)
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", ErrRaffles.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrRaffles, w)
		return
	}
}

// GetByID is an endpoint for getting raffle by id.
// @Summary	Provides raffle by id
// @Tags	Raffles
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Success	200		{object}	RaffleView
// @Failure 400,401,404,500	{object}	common.ErrResponseCode
// @Router	/raffles/{id}/	[get].
func (controller *Raffles) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raffleID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse id")).Serve(controller.log, ErrRaffles, w)
		return
	}

	raffle, err := controller.raffles.Get(ctx, raffleID)
	if err != nil {
		controller.log.Error("failed to get raffle by id", ErrRaffles.Wrap(err))
		if errors.Is(err, raffles.ErrNoRaffle) {
			common.NewErrResponse(http.StatusNotFound, raffles.ErrNoRaffle).Serve(controller.log, ErrRaffles, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get raffle by id")).Serve(controller.log, ErrRaffles, w)
		return
	}

	gifts, err := controller.raffles.ListGifts(ctx, raffle.FundraiseID)
	if err != nil {
		controller.log.Error("error while listing gifts", ErrRaffles.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list raffle gifts")).Serve(controller.log, ErrRaffles, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToRaffleView(raffle, gifts)); err != nil {
		controller.log.Error("error while encoding response", ErrRaffles.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrRaffles, w)
		return
	}
}
