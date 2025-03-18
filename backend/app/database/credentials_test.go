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
	"one-help/app/users/auth"
)

func TestUserCredentials(t *testing.T) {
	user := users.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
	}
	user2 := users.User{
		ID:        uuid.New(),
		FirstName: "Jane",
		LastName:  "Doe",
	}
	creds := auth.Credentials{
		UserID:       user.ID,
		Email:        "john.doe@example.com",
		PhoneNumber:  "1234567890",
		PasswordHash: "hashhashhash",
	}
	creds2 := auth.Credentials{
		UserID:       user2.ID,
		Email:        "jane.doe@example.com",
		PhoneNumber:  "0987654321",
		PasswordHash: "kisskisskiss",
	}

	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		credsRepository := db.Credentials()
		usersRepository := db.Users()

		t.Run("seed", func(t *testing.T) {
			require.NoError(t, usersRepository.Create(ctx, user))
			require.NoError(t, usersRepository.Create(ctx, user2))
		})

		t.Run("Create&Get", func(t *testing.T) {
			err := credsRepository.Create(ctx, creds)
			require.NoError(t, err)

			storedCreds, err := credsRepository.Get(ctx, auth.NewGetByID(creds.UserID))
			require.NoError(t, err)
			assert.Equal(t, creds, storedCreds)
		})

		t.Run("Get(negative)", func(t *testing.T) {
			_, err := credsRepository.Get(ctx, auth.NewGetByID(uuid.New()))
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrNoUserCredentials)
		})

		t.Run("Get by Email", func(t *testing.T) {
			storedCreds, err := credsRepository.Get(ctx, auth.NewGetByEmail(creds.Email))
			require.NoError(t, err)
			assert.Equal(t, creds, storedCreds)
		})

		t.Run("Get by PhoneNumber", func(t *testing.T) {
			storedCreds, err := credsRepository.Get(ctx, auth.NewGetByPhoneNumber(creds.PhoneNumber))
			require.NoError(t, err)
			assert.Equal(t, creds, storedCreds)
		})

		t.Run("Create with existing user credentials", func(t *testing.T) {
			err := credsRepository.Create(ctx, creds)
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrUserAlreadyHasCredentials)
		})

		t.Run("Create with existing email", func(t *testing.T) {
			existingCreds := creds2
			existingCreds.Email = creds.Email

			err := credsRepository.Create(ctx, existingCreds)
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrUserEmailTaken)
		})

		t.Run("Create with existing phone number", func(t *testing.T) {
			existingCreds := creds2
			existingCreds.PhoneNumber = creds.PhoneNumber

			err := credsRepository.Create(ctx, existingCreds)
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrUserPhoneNumberTaken)
		})

		t.Run("Create second creds", func(t *testing.T) {
			require.NoError(t, credsRepository.Create(ctx, creds2))
		})

		t.Run("Update with existing email", func(t *testing.T) {
			existingCreds := creds2
			existingCreds.Email = creds.Email

			err := credsRepository.Update(ctx, existingCreds)
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrUserEmailTaken)
		})

		t.Run("Update with existing phone number", func(t *testing.T) {
			existingCreds := creds2
			existingCreds.PhoneNumber = creds.PhoneNumber

			err := credsRepository.Update(ctx, existingCreds)
			require.Error(t, err)
			require.ErrorIs(t, err, auth.ErrUserPhoneNumberTaken)
		})

		t.Run("Update", func(t *testing.T) {
			creds.Email = "john.doe@newemail.com"
			creds.PhoneNumber = "5555555555"

			storedCreds, err := credsRepository.Get(ctx, auth.NewGetByID(creds.UserID))
			require.NoError(t, err)
			assert.NotEqual(t, creds, storedCreds)

			err = credsRepository.Update(ctx, creds)
			require.NoError(t, err)

			storedCreds, err = credsRepository.Get(ctx, auth.NewGetByID(creds.UserID))
			require.NoError(t, err)
			assert.Equal(t, creds, storedCreds)
		})
	})
}
