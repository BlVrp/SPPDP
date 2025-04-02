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
	"one-help/app/payments"
	"one-help/app/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayments(t *testing.T) {
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

	paymentType := "Stripe"

	payment := payments.Payment{
		DonationId:    donation.ID,
		PaymentType:   paymentType,
		TransactionId: "123456",
		Confirmed:     false,
	}

	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		fundraiseRepository := db.Fundraises()
		usersRepository := db.Users()
		fundraiseStatusesRepository := db.FundraiseStatuses()
		donationsRepository := db.Donations()
		paymentTypesRepository := db.PaymentTypes()
		paymentsRepository := db.Payments()
		t.Run("Create&Get", func(t *testing.T) {
			require.NoError(t, fundraiseStatusesRepository.Create(ctx, fundraiseStatus))
			require.NoError(t, usersRepository.Create(ctx, user))
			require.NoError(t, fundraiseRepository.Create(ctx, fundraise))
			require.NoError(t, donationsRepository.Create(ctx, donation))
			require.NoError(t, paymentTypesRepository.Create(ctx, paymentType))
			require.NoError(t, paymentsRepository.Create(ctx, payment))

			storedPayment, err := paymentsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			assert.Equal(t, payment, storedPayment)
		})

		t.Run("Get(negative)", func(t *testing.T) {
			donationId := uuid.New()
			_, err := paymentsRepository.Get(ctx, donationId)
			require.Error(t, err)
			require.ErrorIs(t, err, payments.ErrNoPayment)
		})

		t.Run("Update", func(t *testing.T) {
			payment.Confirmed = true

			storedPayment, err := paymentsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			assert.NotEqual(t, payment, storedPayment)

			err = paymentsRepository.Update(ctx, payment)
			require.NoError(t, err)

			storedPayment, err = paymentsRepository.Get(ctx, donation.ID)
			require.NoError(t, err)
			assert.Equal(t, payment, storedPayment)
		})

		t.Run("List", func(t *testing.T) {
			storedPayments, err := paymentsRepository.List(ctx)
			require.NoError(t, err)
			assert.Equal(t, payment, storedPayments[0])
			assert.Equal(t, 1, len(storedPayments))
		})

		t.Run("Delete", func(t *testing.T) {
			err := paymentsRepository.Delete(ctx, donation.ID)
			require.NoError(t, err)

			_, err = paymentsRepository.Get(ctx, donation.ID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, payments.ErrNoPayment)
		})
	})
}
