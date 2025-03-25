package database_test

import (
	"context"
	"one-help/app"
	"one-help/app/database"
	"one-help/app/database/dbtesting"
	"one-help/app/fundraises"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFundraises(t *testing.T) {
	fundraise := fundraises.Fundraise{
		ID:           uuid.New(),
		Title:        "Test",
		Description:  "Test Description",
		TargetAmount: 234.4,
		StartDate:    time.Now(),
		EndDate:      time.Time{},
		Status:       "Active",
	}
	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		fundraiseRepository := db.Fundraises()
		t.Run("Create&Get", func(t *testing.T) {
			err := fundraiseRepository.Create(ctx, fundraise)
			require.NoError(t, err)

			storedFundraise, err := fundraiseRepository.Get(ctx, fundraise.ID)
			require.NoError(t, err)
			assert.Equal(t, fundraise, storedFundraise)
		})

		t.Run("Get(negative)", func(t *testing.T) {
			fundraiseId := uuid.New()
			_, err := fundraiseRepository.Get(ctx, fundraiseId)
			require.Error(t, err)
			require.ErrorIs(t, err, fundraises.ErrNoFundraise)
		})

		t.Run("Update", func(t *testing.T) {
			fundraise.Title = "Test 2"
			fundraise.EndDate = time.Now().Add(time.Hour)

			storedFundraise, err := fundraiseRepository.Get(ctx, fundraise.ID)
			require.NoError(t, err)
			assert.NotEqual(t, fundraise, storedFundraise)

			err = fundraiseRepository.Update(ctx, fundraise)
			require.NoError(t, err)

			storedFundraise, err = fundraiseRepository.Get(ctx, fundraise.ID)
			require.NoError(t, err)
			assert.Equal(t, fundraise, storedFundraise)
		})

		t.Run("List", func(t *testing.T) {
			storedFundraises, err := fundraiseRepository.List(ctx)
			require.NoError(t, err)
			assert.Equal(t, fundraise, storedFundraises[0])
			assert.Equal(t, 1, len(storedFundraises))
		})

		t.Run("Delete", func(t *testing.T) {
			err := fundraiseRepository.Delete(ctx, fundraise.ID)
			require.NoError(t, err)

			_, err = fundraiseRepository.Get(ctx, fundraise.ID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fundraises.ErrNoFundraise)
		})
	})
}
