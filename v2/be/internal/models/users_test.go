package models

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestUsersModelCreate(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()
	hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
	require.NoError(t, err)

	username := fmt.Sprintf("testuser_%s", time.Now().String())
	usersModel := &UsersModel{pool: pool}

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "Valid user creation",
			user: &User{
				ID:       uuid.New().String(),
				Username: username,
				Password: []byte(hash),
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &User{
				ID:       uuid.New().String(),
				Username: username,
				Password: []byte(hash),
			},
			wantErr: true,
		},
		{
			name: "Empty username",
			user: &User{
				ID:       uuid.New().String(),
				Username: "",
				Password: []byte(hash),
			},
			wantErr: true,
		},
		{
			name: "Empty password",
			user: &User{
				ID:       uuid.New().String(),
				Username: fmt.Sprintf("testuser_%s", time.Now().String()),
				Password: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usersModel.Create(context.Background(), tt.user)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsersModelCreateWithCancelledContext(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
	require.NoError(t, err)
	usersModel := &UsersModel{pool: pool}

	user := &User{
		ID:       uuid.New().String(),
		Username: "testuser",
		Password: []byte(hash),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = usersModel.Create(ctx, user)
	require.Error(t, err)
}

func TestGetByUsername(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel := &UsersModel{pool: pool}

	t.Run("ExistingUser", func(t *testing.T) {
		ctx := context.Background()
		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}

		err = usersModel.Create(context.Background(), user)
		require.NoError(t, err)

		user, err = usersModel.GetByUsername(ctx, user.Username)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Password)
	})

	t.Run("NonExistentUser", func(t *testing.T) {
		ctx := context.Background()
		username := "nonexistentuser"

		user, err := usersModel.GetByUsername(ctx, username)
		require.Error(t, err)
		require.Nil(t, user)
		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("EmptyUsername", func(t *testing.T) {
		ctx := context.Background()
		username := ""

		user, err := usersModel.GetByUsername(ctx, username)
		require.Error(t, err)
		require.Nil(t, user)
	})

	t.Run("CanceledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		username := "existinguser"

		user, err := usersModel.GetByUsername(ctx, username)
		require.Error(t, err)
		require.Nil(t, user)
	})

	t.Run("TransactionRollback", func(t *testing.T) {
		ctx := context.Background()
		username := "existinguser"

		// Simulate a transaction rollback by closing the pool
		pool.Close()

		user, err := usersModel.GetByUsername(ctx, username)
		require.Error(t, err)
		require.Nil(t, user)
	})
}

func TestUsersModelExists(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel := &UsersModel{pool: pool}

	t.Run("Existing User", func(t *testing.T) {
		ctx := context.Background()
		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}

		err = usersModel.Create(context.Background(), user)
		require.NoError(t, err)

		exists, err := usersModel.Exists(ctx, user.ID)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("Non-Existent User", func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New().String()

		exists, err := usersModel.Exists(ctx, id)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		ctx := context.Background()
		id := "invalid-id-format"

		exists, err := usersModel.Exists(ctx, id)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("Empty ID", func(t *testing.T) {
		ctx := context.Background()
		exists, err := usersModel.Exists(ctx, "")
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func TestUsersModelExistsWithClosedPool(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	usersModel := &UsersModel{pool: pool}

	pool.Close()

	ctx := context.Background()
	_, err = usersModel.Exists(ctx, "some-id")
	require.Error(t, err)
}

func TestUsersModelExistsWithCancelledContext(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel := &UsersModel{pool: pool}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = usersModel.Exists(ctx, "some-id")
	require.Error(t, err)
}
