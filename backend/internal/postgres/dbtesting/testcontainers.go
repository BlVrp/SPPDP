package dbtesting

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	// defaultPostgresImage defines default postgres docker image.
	defaultPostgresImage = "postgres:alpine"
)

var (
	// containerMap holds running containers using connection string as key.
	containerMap sync.Map
)

// SetupPgContainer runs and returns postgres test container.
// The container will be terminated by testcontainers.Reaper when all tests are finished.
func SetupPgContainer(ctx context.Context, connStr string) (*testpostgres.PostgresContainer, error) {
	if container, exist := containerMap.Load(connStr); exist {
		return container.(*testpostgres.PostgresContainer), nil
	}

	connParams, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	pqPass, ok := connParams.User.Password()
	if !ok {
		return nil, errors.New("couldn't get pq password from conn string")
	}

	pgContainer, err := testpostgres.Run(
		ctx,
		defaultPostgresImage,
		testpostgres.WithDatabase(strings.Trim(connParams.Path, "/")),
		testpostgres.WithUsername(connParams.User.Username()),
		testpostgres.WithPassword(pqPass),
		testpostgres.WithSQLDriver(connParams.Scheme),
		testpostgres.BasicWaitStrategies(),
		testcontainers.WithLogger(testcontainers.Logger),
	)
	if err != nil {
		return nil, err
	}

	containerMap.Store(connStr, pgContainer)

	return pgContainer, nil
}
