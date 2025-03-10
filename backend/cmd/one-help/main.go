package main

import (
	"context"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	onehelp "one-help"
	"one-help/app/database"
	"one-help/internal/logger/zaplog"
	"one-help/internal/process"
)

// Error is a default error type for one-help cli.
var Error = errs.Class("one-help cli")

// Config is the global configuration for one-help.
type Config struct {
	Database   database.Config
	Config     onehelp.Config
	ServerName string `env:"SERVER_NAME"`
}

// commands.
var (
	rootCmd = &cobra.Command{
		Use:   "one-help",
		Short: "cli for interacting with one-help service",
	}
	runCmd = &cobra.Command{
		Use:         "run",
		Short:       "runs the program",
		RunE:        cmdRun,
		Annotations: map[string]string{"type": "run"},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	process.OnSigInt(func() {
		// starting graceful exit on context cancellation.
		cancel()
	})

	log := zaplog.NewLog()

	if err = godotenv.Overload("./configs/.one-help.env"); err != nil {
		log.Error("could not load launchpad config: %v", Error.Wrap(err))
		return Error.Wrap(err)
	}

	config := new(Config)
	envOpt := env.Options{RequiredIfNoDef: true}
	if err = env.Parse(config, envOpt); err != nil {
		log.Error("could not parse config: %v", Error.Wrap(err))
		return Error.Wrap(err)
	}

	db, err := database.New(config.Database.Database)
	if err != nil {
		log.Error("error connecting to launchpad database", Error.Wrap(err))
		return Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	if err = db.ExecuteMigrations(config.Database.MigrationsPath, true); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			log.Error("error migration launchpad db", Error.Wrap(err))
			return Error.Wrap(err)
		}
	}

	peer, err := onehelp.New(ctx, log, config.Config, db)
	if err != nil {
		log.Error("could not initialize peer", Error.Wrap(err))
		return Error.Wrap(err)
	}

	return errs.Combine(peer.Run(ctx), peer.Close())
}
