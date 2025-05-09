//nolint:godot,revive
package console

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpswagger "github.com/swaggo/http-swagger"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"one-help/app/console/controllers/common"
	eventscontroller "one-help/app/console/controllers/events"
	fundraisescontroller "one-help/app/console/controllers/fundraises"
	infocontroller "one-help/app/console/controllers/info"
	rafflescontroller "one-help/app/console/controllers/raffles"
	userscontroller "one-help/app/console/controllers/users"
	_ "one-help/app/console/docs"
	"one-help/app/events"
	"one-help/app/fundraises"
	"one-help/app/raffles"
	"one-help/app/users"
	"one-help/internal/logger"
)

const (
	bearerPrefix = "Bearer "
)

var (
	// Error is an error class that indicates internal http server error.
	Error = errs.Class("console web server error")
)

// Config contains configuration for console web server.
type Config struct {
	Address                    string `env:"ADDRESS"`
	FrontEndPaymentRedirectUrl string `env:"FRONT_END_PAYMENT_REDIRECT_URL"`
}

// Server represents console web server.
//
// architecture: Endpoint.
type Server struct {
	log    logger.Logger
	config Config

	listener net.Listener
	server   http.Server

	users      *users.Service
	fundraises *fundraises.Service
	events     *events.Service
	raffles    *raffles.Service
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
	users *users.Service,
	fundraises *fundraises.Service,
	events *events.Service,
	raffles *raffles.Service,
) *Server {
	server := &Server{
		log:        log,
		config:     config,
		listener:   listener,
		users:      users,
		fundraises: fundraises,
		events:     events,
		raffles:    raffles,
	}

	infoController := infocontroller.NewInfo(log)
	usersController := userscontroller.NewUsers(log, users)
	fundraisesController := fundraisescontroller.NewFundraises(log, fundraises, config.FrontEndPaymentRedirectUrl)
	eventsController := eventscontroller.NewEvents(log, events, fundraises)
	rafflesController := rafflescontroller.NewRaffles(log, raffles, fundraises)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v0").Subrouter()

	infoRouter := apiRouter.PathPrefix("/info").Subrouter()
	infoRouter.Use(server.jsonResponse)
	infoRouter.StrictSlash(true)
	infoRouter.HandleFunc("/messages", infoController.Messages).Methods(http.MethodGet, http.MethodOptions)

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.Use(server.jsonResponse)
	authRouter.StrictSlash(true)
	authRouter.HandleFunc("/login", usersController.Login).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/register", usersController.Register).Methods(http.MethodPost, http.MethodOptions)

	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	usersRouter.Use(server.jsonResponse)
	usersRouter.Use(server.withAuthMiddleware)
	usersRouter.StrictSlash(true)
	usersRouter.HandleFunc("/", usersController.Get).Methods(http.MethodGet, http.MethodOptions)
	usersRouter.HandleFunc("/", usersController.Update).Methods(http.MethodPatch, http.MethodOptions)
	usersRouter.HandleFunc("/change-password", usersController.ChangePassword).Methods(http.MethodPatch, http.MethodOptions)
	usersRouter.HandleFunc("/{id}", usersController.GetByID).Methods(http.MethodGet, http.MethodOptions)
	usersRouter.HandleFunc("/raffle-participants/{id}", usersController.GetRaffleParticipants).Methods(http.MethodGet, http.MethodOptions)

	fundraisesRouter := apiRouter.PathPrefix("/fundraises").Subrouter()
	fundraisesRouter.Use(server.jsonResponse)
	fundraisesRouter.Use(server.withAuthMiddleware)
	fundraisesRouter.StrictSlash(true)
	fundraisesRouter.HandleFunc("/", fundraisesController.List).Methods(http.MethodGet, http.MethodOptions)
	fundraisesRouter.HandleFunc("/my", fundraisesController.ListMy).Methods(http.MethodGet, http.MethodOptions)
	fundraisesRouter.HandleFunc("/{id}", fundraisesController.GetByID).Methods(http.MethodGet, http.MethodOptions)
	fundraisesRouter.HandleFunc("/", fundraisesController.Create).Methods(http.MethodPost, http.MethodOptions)
	fundraisesRouter.HandleFunc("/{id}/donate", fundraisesController.Donate).Methods(http.MethodPost, http.MethodOptions)

	donationsRouter := apiRouter.PathPrefix("/fundraises/donations").Subrouter()
	donationsRouter.Use(server.jsonResponse)
	donationsRouter.StrictSlash(true)
	donationsRouter.HandleFunc("/{id}", fundraisesController.FinishDonation).Methods(http.MethodGet, http.MethodOptions)

	eventsRouter := apiRouter.PathPrefix("/events").Subrouter()
	eventsRouter.Use(server.jsonResponse)
	eventsRouter.Use(server.withAuthMiddleware)
	eventsRouter.StrictSlash(true)
	eventsRouter.HandleFunc("/", eventsController.List).Methods(http.MethodGet, http.MethodOptions)
	eventsRouter.HandleFunc("/{id}", eventsController.GetByID).Methods(http.MethodGet, http.MethodOptions)
	eventsRouter.HandleFunc("/", eventsController.Create).Methods(http.MethodPost, http.MethodOptions)
	eventsRouter.HandleFunc("/{id}/enroll", eventsController.Enroll).Methods(http.MethodPost, http.MethodOptions)

	rafflesRouter := apiRouter.PathPrefix("/raffles").Subrouter()
	rafflesRouter.Use(server.jsonResponse)
	rafflesRouter.Use(server.withAuthMiddleware)
	rafflesRouter.StrictSlash(true)
	rafflesRouter.HandleFunc("/", rafflesController.List).Methods(http.MethodGet, http.MethodOptions)
	rafflesRouter.HandleFunc("/{id}", rafflesController.GetByID).Methods(http.MethodGet, http.MethodOptions)
	rafflesRouter.HandleFunc("/", rafflesController.Create).Methods(http.MethodPost, http.MethodOptions)

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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		handler.ServeHTTP(w, r.Clone(r.Context()))
	})
}

// withAuth performs initial authorization before every request.
func (server *Server) withAuth(handler http.Handler, optionalAuth bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("Authorization")
		if token == "" {
			if optionalAuth {
				handler.ServeHTTP(w, r.Clone(ctx))
				return
			}

			common.NewErrResponse(http.StatusUnauthorized, errors.New("token is empty")).Serve(server.log, Error, w)
			return
		}

		token = strings.TrimPrefix(token, bearerPrefix)

		jwt, err := server.users.ValidateJWTToken(ctx, token)
		if err != nil {
			common.NewErrResponse(http.StatusUnauthorized, errors.Unwrap(err)).Serve(server.log, Error, w)
			return
		}

		ctx = jwt.Payload.SetIntoContext(ctx)

		handler.ServeHTTP(w, r.Clone(ctx))
	})
}

func (server *Server) withAuthMiddleware(handler http.Handler) http.Handler {
	return server.withAuth(handler, false)
}
