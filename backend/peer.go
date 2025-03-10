package backend

import (
	"context"
	"net"

	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"one-help/app"
	"one-help/app/console"
	"one-help/internal/logger"
)

// Config is the global configuration for one-help app.
type Config struct {
	Console struct {
		Config console.Config
	}
}

// Peer is the representation of a server.
type Peer struct {
	Config   Config
	Log      logger.Logger
	Database app.DB

	// Console web server with web UI.
	Console struct {
		Listener net.Listener
		Endpoint *console.Server
	}
}

// New is a constructor for peer.
func New(ctx context.Context, logger logger.Logger, config Config, db app.DB) (peer *Peer, err error) {
	peer = &Peer{
		Log:      logger,
		Database: db,
		Config:   config,
	}

	// console setup
	{
		peer.Console.Listener, err = net.Listen("tcp", config.Console.Config.Address)
		if err != nil {
			return &Peer{}, err
		}

		peer.Console.Endpoint = console.NewServer(
			config.Console.Config,
			peer.Log,
			peer.Console.Listener,
		)
	}

	return peer, nil
}

// Run runs console until it's either closed or it errors.
func (peer *Peer) Run(ctx context.Context) error {
	peer.Log.Info("one-help running")

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return peer.Console.Endpoint.Run(ctx)
	})

	return group.Wait()
}

// Close closes all the resources.
func (peer *Peer) Close() error {
	peer.Log.Debug("one-help closing")
	var errlist errs.Group

	if peer.Console.Endpoint != nil {
		errlist.Add(peer.Console.Endpoint.Close())
	}

	return nil
}
