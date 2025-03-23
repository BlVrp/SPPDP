package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"one-help/app"
	"one-help/app/database"
	"one-help/app/database/dbtesting"
)

func TestFundraiseStatuses(t *testing.T) {
	status := "Status 1"
	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		fundraiseStatusesRepository := db.FundraiseStatuses()
		t.Run("Create&List", func(t *testing.T) {
			err := fundraiseStatusesRepository.Create(ctx, status)
			require.NoError(t, err)
			storedStatuses, err := fundraiseStatusesRepository.List(ctx)
			require.NoError(t, err)
			assert.Equal(t, status, storedStatuses[0])
		})

		t.Run("Delete", func(t *testing.T) {
			err := fundraiseStatusesRepository.Delete(ctx, status)
			require.NoError(t, err)
			storedStatuses, err := fundraiseStatusesRepository.List(ctx)
			require.NoError(t, err)
			assert.Empty(t, storedStatuses)
		})
	})
}
