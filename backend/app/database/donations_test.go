package database_test

import (
	"context"
	"testing"
	"time"

	"one-help/app"
	"one-help/app/database"
	"one-help/app/database/dbtesting"
	"one-help/app/donations"
	"one-help/app/fundraises"
	"one-help/app/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDonations(t *testing.T) {
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

	donation := donations.Donation{
		ID:          uuid.New(),
		UserId:      user.ID,
		FundraiseId: fundraise.ID,
		Amount:      100.0,
		CreatedAt:   time.Now(),
	}

	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		fundraiseRepository := db.Fundraises()
		usersRepository := db.Users()
		fundraiseStatusesRepository := db.FundraiseStatuses()
		donationsRepository := db.Donations()
		t.Run("Create&Get", func(t *testing.T) {
			require.NoError(t, fundraiseStatusesRepository.Create(ctx, fundraiseStatus))
			require.NoError(t, usersRepository.Create(ctx, user))
			require.NoError(t, fundraiseRepository.Create(ctx, fundraise))
			require.NoError(t, donationsRepository.Create(ctx, donation))

			storedDonation, err := donationsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			donationsAreEqual(t, donation, storedDonation)
		})

		t.Run("Get(negative)", func(t *testing.T) {
			donationId := uuid.New()
			_, err := donationsRepository.Get(ctx, donationId)
			require.Error(t, err)
			require.ErrorIs(t, err, donations.ErrNoDonation)
		})

		t.Run("Update", func(t *testing.T) {
			donation.Amount = 25

			storedDonation, err := donationsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			assert.NotEqual(t, donation, storedDonation)

			err = donationsRepository.Update(ctx, donation)
			require.NoError(t, err)

			storedDonation, err = donationsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			donationsAreEqual(t, donation, storedDonation)
		})

		t.Run("List", func(t *testing.T) {
			storedDonations, err := donationsRepository.List(ctx, donations.ListParams{})
			require.NoError(t, err)
			donationsAreEqual(t, donation, storedDonations[0])
			assert.Equal(t, 1, len(storedDonations))
		})

		t.Run("Delete", func(t *testing.T) {
			err := donationsRepository.Delete(ctx, donation.ID)
			require.NoError(t, err)

			_, err = donationsRepository.Get(ctx, donation.ID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, donations.ErrNoDonation)
		})
	})
}

func donationsAreEqual(t *testing.T, expected, actual donations.Donation) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.FundraiseId, actual.FundraiseId)
	assert.Equal(t, expected.UserId, actual.UserId)
	assert.Equal(t, expected.Amount, actual.Amount)
}
