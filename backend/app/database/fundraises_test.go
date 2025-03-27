package database_test

import (
	"context"
	"testing"
	"time"

	"one-help/app"
	"one-help/app/database"
	"one-help/app/database/dbtesting"
	"one-help/app/fundraises"
	"one-help/app/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFundraises(t *testing.T) {
	user := users.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Website:   "https://example.com",
		FileName:  "john_doe.txt",
	}

	fundraiseStatus := "Active"

	fundraise := fundraises.Fundraise{
		ID:           uuid.New(),
		OrganizerId:  user.ID,
		Title:        "Test",
		Description:  "Test Description",
		TargetAmount: 234.4,
		StartDate:    time.Now(),
		EndDate:      time.Time{},
		Status:       fundraiseStatus,
	}

	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		fundraiseRepository := db.Fundraises()
		usersRepository := db.Users()
		fundraiseStatusesRepository := db.FundraiseStatuses()
		t.Run("Create&Get", func(t *testing.T) {
			require.NoError(t, fundraiseStatusesRepository.Create(ctx, fundraiseStatus))
			require.NoError(t, usersRepository.Create(ctx, user))
			require.NoError(t, fundraiseRepository.Create(ctx, fundraise))

			storedFundraise, err := fundraiseRepository.Get(ctx, fundraise.ID)
			require.NoError(t, err)
			fundraisesAreEqual(t, fundraise, storedFundraise)
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
			fundraisesAreEqual(t, fundraise, storedFundraise)
		})

		t.Run("List", func(t *testing.T) {
			storedFundraises, err := fundraiseRepository.List(ctx)
			require.NoError(t, err)
			fundraisesAreEqual(t, fundraise, storedFundraises[0])
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

func fundraisesAreEqual(t *testing.T, expected, actual fundraises.Fundraise) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.OrganizerId, actual.OrganizerId)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.TargetAmount, actual.TargetAmount)
	assert.Equal(t, expected.Status, actual.Status)
}
