package database_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"one-help/app"
	"one-help/app/database"
	"one-help/app/database/dbtesting"
	"one-help/app/users"
)

func TestUsers(t *testing.T) {
	user := users.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Website:   "https://example.com",
		ImageUrl:  "john_doe.txt",
	}
	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		usersRepository := db.Users()
		t.Run("Create&Get", func(t *testing.T) {
			err := usersRepository.Create(ctx, user)
			require.NoError(t, err)

			storedUser, err := usersRepository.Get(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user, storedUser)
		})

		t.Run("Get(negative)", func(t *testing.T) {
			userID := uuid.New()
			_, err := usersRepository.Get(ctx, userID)
			require.Error(t, err)
			require.ErrorIs(t, err, users.ErrNoUser)
		})

		t.Run("Update", func(t *testing.T) {
			user.EmptyDeliveryAddress()
			user.LastName = "Wog"

			storedUser, err := usersRepository.Get(ctx, user.ID)
			require.NoError(t, err)
			assert.NotEqual(t, user, storedUser)

			err = usersRepository.Update(ctx, user)
			require.NoError(t, err)

			storedUser, err = usersRepository.Get(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user, storedUser)
		})

		t.Run("Delete", func(t *testing.T) {
			err := usersRepository.Delete(ctx, user.ID)
			require.NoError(t, err)

			_, err = usersRepository.Get(ctx, user.ID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, users.ErrNoUser)
		})
	})
}
