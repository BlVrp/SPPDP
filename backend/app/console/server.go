//nolint:godot,revive
package console

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	httpswagger "github.com/swaggo/http-swagger"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	infocontroller "one-help/app/console/controllers/info"
	_ "one-help/app/console/docs"
	"one-help/internal/logger"
)

var (
	// Error is an error class that indicates internal http server error.
	Error = errs.Class("console web server error")
)

// Config contains configuration for console web server.
type Config struct {
	Address string `env:"ADDRESS"`
}

// Server represents console web server.
//
// architecture: Endpoint.
type Server struct {
	log    logger.Logger
	config Config

	listener net.Listener
	server   http.Server
}

// NewServer is a constructor for console web server.
// @title		Swagger API for One-Help app.
// @version		v0
// @description	This is Swagger's auto-generated documentation for One-Help app.
// @BasePath  /api/v0
func NewServer(
	config Config,
	log logger.Logger,
	listener net.Listener,
) *Server {
	server := &Server{
		log:      log,
		config:   config,
		listener: listener,
	}

	infoController := infocontroller.NewInfo(log)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v0").Subrouter()

	infoRouter := apiRouter.PathPrefix("/info").Subrouter()
	infoRouter.Use(server.jsonResponse)
	infoRouter.StrictSlash(true)
	infoRouter.HandleFunc("/messages", infoController.Messages).Methods(http.MethodGet, http.MethodOptions)

	apiRouter.PathPrefix("/docs/swagger/").Handler(httpswagger.WrapHandler)

	server.server = http.Server{
		Handler: router,
	}

	return server
}

// Run starts the server that host webapp and api endpoint.
func (server *Server) Run(ctx context.Context) (err error) {
	var group errgroup.Group

	ctx, cancel := context.WithCancel(ctx)

	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.server.Shutdown(ctx))
	})

	group.Go(func() error {
		defer cancel()

		server.log.InfoF("server is running on: http://%s", server.config.Address)
		err = server.server.Serve(server.listener)
		isCancelled := errs.IsFunc(err, func(err error) bool { return errors.Is(err, context.Canceled) })
		if isCancelled || errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return Error.Wrap(err)
	})

	return Error.Wrap(group.Wait())
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.server.Close())
}

// jsonResponse sets a response "Content-Type" value as "application/json"
func (server *Server) jsonResponse(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r.Clone(r.Context()))
	})
}
