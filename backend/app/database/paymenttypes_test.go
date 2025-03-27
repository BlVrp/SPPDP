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

func TestPaymentTypes(t *testing.T) {
	paymentType := "Type 1"
	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		paymentTypesRepository := db.PaymentTypes()
		t.Run("Create&List", func(t *testing.T) {
			err := paymentTypesRepository.Create(ctx, paymentType)
			require.NoError(t, err)
			storedTypes, err := paymentTypesRepository.List(ctx)
			require.NoError(t, err)
			assert.Equal(t, paymentType, storedTypes[0])
			assert.Equal(t, 1, len(storedTypes))
		})

		t.Run("Delete", func(t *testing.T) {
			err := paymentTypesRepository.Delete(ctx, paymentType)
			require.NoError(t, err)
			storedTypes, err := paymentTypesRepository.List(ctx)
			require.NoError(t, err)
			assert.Empty(t, storedTypes)
		})
	})
}
