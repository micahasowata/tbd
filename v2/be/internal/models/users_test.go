package models_test

import (
	"context"
	"testing"

	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestUsersCreate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)
	})

	t.Run("duplicate username", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		uTwo := &models.User{
			ID:       db.NewID(),
			Username: u.Username,
			Password: []byte(testUserPassword(t)),
		}

		err = users.Create(context.Background(), uTwo)
		require.ErrorIs(t, err, models.ErrDuplicateUsername)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := users.Create(ctx, u)
		require.Error(t, err)
	})
}

func TestGetByUsername(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		user, err := users.GetByUsername(context.Background(), u.Username)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, u.ID, user.ID)
		require.Equal(t, u.Password, user.Password)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}

		u, err := users.GetByUsername(context.Background(), gofakeit.AdverbManner())
		require.Error(t, err)
		require.Nil(t, u)
		require.ErrorIs(t, err, models.ErrRecordNotFound)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		u, err := users.GetByUsername(ctx, gofakeit.Username())
		require.Error(t, err)
		require.Nil(t, u)
	})
}

func TestUsersModelExists(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		exists, err := users.Exists(context.Background(), u.ID)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}

		exists, err := users.Exists(context.Background(), db.NewID())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)
		users := &models.UsersModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		exists, err := users.Exists(ctx, db.NewID())
		require.Error(t, err)
		require.False(t, exists)
	})
}
