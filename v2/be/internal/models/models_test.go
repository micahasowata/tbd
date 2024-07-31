package models

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	models := New(pool)
	require.NotNil(t, models)
	require.NotNil(t, models.Users)
	require.NotNil(t, models.Tasks)
	require.NotNil(t, models.Users.pool)
	require.NotNil(t, models.Tasks.pool)
	require.Equal(t, models.Users.pool, pool)
	require.Equal(t, models.Tasks.pool, pool)
}

func TestNewWithNilPool(t *testing.T) {
	models := New(nil)

	require.NotNil(t, models)
	require.NotNil(t, models.Users)
	require.Nil(t, models.Users.pool)
	require.NotNil(t, models.Tasks)
	require.Nil(t, models.Tasks.pool)
}
