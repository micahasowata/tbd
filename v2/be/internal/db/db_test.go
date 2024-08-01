package db_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"v2/be/internal/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		// setup dockertest
		pool, err := dockertest.NewPool("")
		require.Nil(t, err)

		// Ensure docker engine is running
		err = pool.Client.Ping()
		require.Nil(t, err)

		// Run container with option to make it instantly removable
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16.3-alpine3.20",
			Env: []string{
				"POSTGRES_PASSWORD=pa55word",
				"POSTGRES_USER=tester",
				"POSTGRES_DB=tdb",
			},
		}, func(c *docker.HostConfig) {
			c.AutoRemove = true
			c.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

		require.Nil(t, err)

		dsn := fmt.Sprintf("postgresql://tester:pa55word@localhost:%s/tdb?sslmode=disable", resource.GetPort("5432/tcp"))

		// Wait till container is ready for connections
		err = pool.Retry(func() error {
			p, err := pgxpool.New(context.Background(), dsn)
			require.Nil(t, err)

			return p.Ping(context.Background())
		})
		if err != nil {
			t.Log(dsn, err.Error())
		}
		require.Nil(t, err)

		// Test db.New()
		d, err := db.New(dsn)
		require.Nil(t, err)
		require.NotNil(t, d)

		// Remove container
		defer func() {
			err := pool.Purge(resource)
			require.Nil(t, err)
		}()
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name string
			dsn  string
		}{
			{
				name: "no dsn",
				dsn:  "",
			},
			{
				name: "invalid dsn",
				dsn:  "invalid://dsn",
			},
			{
				name: "fake dsn",
				dsn:  "postgresql://users:se[ret@localhost:5432/db?sslmode=disable",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d, err := db.New(tt.dsn)

				require.NotNil(t, err)
				require.Nil(t, d)
			})
		}
	})
}

func TestFormatErr(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect string
	}{
		{
			name:   "valid",
			err:    &pgconn.PgError{ConstraintName: "unique_constraint"},
			expect: "unique_constraint",
		},
		{
			name:   "generic",
			err:    errors.New("unique"),
			expect: "",
		},
		{
			name:   "nil",
			err:    nil,
			expect: "",
		},
		{
			name:   "no constraint name",
			err:    &pgconn.PgError{},
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.FormatErr(tt.err)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestNewID(t *testing.T) {
	id := db.NewID()

	require.NotEmpty(t, id)
	require.Len(t, id, 36)
}
