package events

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"one-help/app/console/controllers/common"
	"one-help/app/events"
	"one-help/app/fundraises"
	"one-help/app/users/credentials"
	"one-help/internal/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"
)

var (
	// ErrEvents is an internal error type for events controller.
	ErrEvents = errs.Class("events controller")
)

// Events is a controller that handles all events related routes.
type Events struct {
	log logger.Logger

	events     *events.Service
	fundraises *fundraises.Service
}

// NewEvents is a constructor for events controller.
func NewEvents(log logger.Logger, events *events.Service, fundraises *fundraises.Service) *Events {
	eventsController := &Events{
		log:        log,
		events:     events,
		fundraises: fundraises,
	}

	return eventsController
}

// Create is an endpoint for creating new event.
// @Summary	Creates new event
// @Tags	Events
// @Accept	json
// @Produce	json
// @Param	request	body	CreateRequest	true	"Event needed data fields"
// @Success	200		{object}	EventsView
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/events/	[post].
func (controller *Events) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.log.Error("failed to decode create request body", ErrEvents.Wrap(err))
		common.NewErrResponse(http.StatusBadRequest, err).Serve(controller.log, ErrEvents, w)
		return
	}

	// INFO: Caller creds.
	creds, err := credentials.GetFromContext(ctx)
	if err != nil {
		common.NewErrResponse(http.StatusUnauthorized, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
		return
	}

	fundraise, err := controller.fundraises.Get(ctx, request.FundraiseId)
	if err != nil {
		controller.log.Error("failed to get the fundraise related to the new event", ErrEvents.Wrap(err))
		if fundraises.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
			return
		}
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
		return
	}

	if fundraise.OrganizerId != creds.UserID {
		controller.log.Error("user is not allowed to create event for this fundraise", ErrEvents.New("user is not allowed to create event for this fundraise"))
		common.NewErrResponse(http.StatusForbidden, errors.New("user is not allowed to create event for this fundraise")).Serve(controller.log, ErrEvents, w)
		return
	}

	createParams := events.CreateParams{
		Title:           request.Title,
		Description:     request.Description,
		StartDate:       request.StartDate,
		EndDate:         request.EndDate,
		Format:          request.Format,
		MaxParticipants: request.MaxParticipants,
		MinimumDonation: request.MinimumDonation,
		Address:         request.Address,
		FundraiseId:     request.FundraiseId,
	}

	event, err := controller.events.Create(ctx, createParams)
	if err != nil {
		controller.log.Error("error while creating event", ErrEvents.Wrap(err))
		if events.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToEventView(event)); err != nil {
		controller.log.Error("error while encoding response:", ErrEvents.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
		return
	}
}

// List is an endpoint for listing events.
// @Summary	Returns list of events
// @Tags	Events
// @Produce	json
// @Param	Authorization	header	string	true	"Bearer token to authorize access"
// @Param	limit			query	integer	false	"Items per page (positive number expected) [default value: 20]"
// @Param	page			query	integer	false	"Number of the page (1...) [default value: 1]"
// @Success	200		{object}	common.Page[EventsView]
// @Failure	400,401,500	{object}	common.ErrResponseCode
// @Router	/events/	[get].
func (controller *Events) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		limit = 20
		page  = 1
		err   error
	)
	if val := r.URL.Query().Get("limit"); val != "" {
		limit, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'limit' query parameter", ErrEvents.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid limit value")).Serve(controller.log, ErrEvents, w)
			return
		}
	}
	if val := r.URL.Query().Get("page"); val != "" {
		page, err = strconv.Atoi(val)
		if err != nil {
			controller.log.Error("failed to parse 'page' query parameter", ErrEvents.Wrap(err))
			common.NewErrResponse(http.StatusBadRequest, errors.New("invalid page value")).Serve(controller.log, ErrEvents, w)
			return
		}
	}

	list, err := controller.events.List(ctx, limit, page, nil)
	if err != nil {
		controller.log.Error("failed to list events", ErrEvents.Wrap(err))
		if events.ParamsError.Has(err) {
			common.NewErrResponse(http.StatusBadRequest, errors.Unwrap(err)).Serve(controller.log, ErrEvents, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to list events")).Serve(controller.log, ErrEvents, w)
		return
	}

	var viewList = make([]EventView, len(list))
	for i, event := range list {
		viewList[i] = ToEventView(&event)
	}

	resp := &common.Page[EventView]{
		Data:  viewList,
		Page:  page,
		Limit: limit,
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", ErrEvents.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrEvents, w)
		return
	}
}

// GetByID is an endpoint for getting event by id.
// @Summary	Provides event by id
// @Tags	Event
// @Produce	json
// @Param	Authorization	header	string	false	"Bearer token to authorize access"
// @Success	200		{object}	EventView
// @Failure 400,401,404,500	{object}	common.ErrResponseCode
// @Router	/events/{id}/	[get].
func (controller *Events) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		common.NewErrResponse(http.StatusBadRequest, errs.New("failed to parse id")).Serve(controller.log, ErrEvents, w)
		return
	}

	event, err := controller.events.Get(ctx, eventID)
	if err != nil {
		controller.log.Error("failed to get event by id", ErrEvents.Wrap(err))
		if errors.Is(err, events.ErrNoEvents) {
			common.NewErrResponse(http.StatusNotFound, events.ErrNoEvents).Serve(controller.log, ErrEvents, w)
			return
		}

		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to get event by id")).Serve(controller.log, ErrEvents, w)
		return
	}

	if err = json.NewEncoder(w).Encode(ToEventView(event)); err != nil {
		controller.log.Error("error while encoding response", ErrEvents.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, err).Serve(controller.log, ErrEvents, w)
		return
	}
}
