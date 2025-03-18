//nolint:godot
package info

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zeebo/errs"

	"one-help/app/console/controllers/common"
	"one-help/internal/logger"
)

var (
	// Error is an internal error type for info controller.
	Error = errs.Class("info controller")
)

// Info is a controller that handles all information-related routes.
type Info struct {
	log logger.Logger
}

// NewInfo is a constructor for info controller.
func NewInfo(log logger.Logger) *Info { return &Info{log: log} }

// Messages defines handler for /messages endpoint.
// @Summary	temp endpoint for testing, provides json object with say field containing 'Hello, World!'.
// @Tags	Info
// @Produce	json
// @Success	200	{object}	MessagesResponse
// @Failure	500	{object}	common.ErrResponseCode
// @Router	/info/messages	[get]
func (controller *Info) Messages(w http.ResponseWriter, r *http.Request) {
	resp := &MessagesResponse{Say: "Hello, World!"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		controller.log.Error("error while encoding response", Error.Wrap(err))
		common.NewErrResponse(http.StatusInternalServerError, errors.New("failed to encode response")).Serve(controller.log, Error, w)
		return
	}
}
