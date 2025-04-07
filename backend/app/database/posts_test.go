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

func TestPosts(t *testing.T) {
	post := "Post 1"
	dbtesting.Run(t, database.Config{}, func(ctx context.Context, t *testing.T, db app.DB) {
		postsRepository := db.Posts()
		t.Run("Create&List", func(t *testing.T) {
			err := postsRepository.Create(ctx, post)
			require.NoError(t, err)
			storedPosts, err := postsRepository.List(ctx)
			require.NoError(t, err)
			require.Equal(t, 2, len(storedPosts))
			assert.Contains(t, storedPosts, post)
		})

		t.Run("Delete", func(t *testing.T) {
			err := postsRepository.Delete(ctx, post)
			require.NoError(t, err)
			storedPosts, err := postsRepository.List(ctx)
			require.NoError(t, err)
			require.Equal(t, 1, len(storedPosts))
		})
	})
}
