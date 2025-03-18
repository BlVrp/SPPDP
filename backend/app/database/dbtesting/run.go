package dbtesting

import (
	"context"
	"net/url"
	"testing"

	"one-help/app"
	"one-help/app/database"
	"one-help/internal/postgres"
	"one-help/internal/postgres/dbtesting"
	"one-help/internal/tempdb"
)

// tempMasterDB is a bridge.DB-implementing type that cleans up after itself when closed.
type tempMasterDB struct {
	app.DB
	tempDB *tempdb.TempDatabase
}

const (
	// defaultTestConn defines default test connection string that is expected to work with postgres server.
	defaultTestConn = "postgres://postgres:123456@localhost/one_help_db?sslmode=disable"
	// defaultMigrationsPath defines default path to migration files.
	defaultMigrationsPath = "./migrations"
)

// Run method will establish connection with db, create tables in random schema, run tests.
func Run(t *testing.T, cfg database.Config, test func(ctx context.Context, t *testing.T, db app.DB)) {
	t.Run("Postgres", func(t *testing.T) {
		ctx := context.Background()

		if cfg.Database == "" {
			cfg.Database = defaultTestConn
		}

		if cfg.MigrationsPath == "" {
			cfg.MigrationsPath = defaultMigrationsPath
		}

		connParams, err := url.Parse(cfg.Database)
		if err != nil {
			t.Fatal(err)
		}

		pgContainer, err := dbtesting.SetupPgContainer(ctx, cfg.Database)
		if err != nil {
			t.Fatal(err)
		}

		containerConnStr, err := pgContainer.ConnectionString(ctx, connParams.RawQuery)
		if err != nil {
			t.Fatal(err)
		}

		schema := postgres.SchemaName(t.Name(), "Test", 0, "")

		tempDB, err := tempdb.OpenUnique(ctx, containerConnStr, schema)
		if err != nil {
			t.Fatal(err)
		}

		masterDB, err := CreateMasterDBOnTopOf(tempDB)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			err = tempDB.Cleanup()
			if err != nil {
				t.Fatal(err)
			}

			err = masterDB.Close()
			if err != nil {
				t.Fatal(err)
			}
		})

		const isUp = true
		err = masterDB.ExecuteMigrations(cfg.MigrationsPath, isUp)
		if err != nil {
			t.Fatal(err)
		}

		test(ctx, t, masterDB)
	})
}

// CreateMasterDBOnTopOf creates a new bridge.DB on top of an already existing temporary DB.
func CreateMasterDBOnTopOf(tempDB *tempdb.TempDatabase) (app.DB, error) {
	masterDB, err := database.New(tempDB.ConnStr)
	return &tempMasterDB{DB: masterDB, tempDB: tempDB}, err
}
