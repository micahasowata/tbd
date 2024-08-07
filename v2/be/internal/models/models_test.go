package models_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// dsn stores connection dsn for the models test package
var dsn string

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := db.New(dsn)
	require.NoError(t, err)
	require.Nil(t, err)

	return pool
}

func testUserPassword(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, false, false, 12)
}

func TestMain(m *testing.M) {
	l := zap.NewNop()

	// Setup container
	pool, err := dockertest.NewPool("")
	if err != nil {
		l.Fatal(err.Error(), zap.Error(err))
	}

	err = pool.Client.Ping()
	if err != nil {
		l.Fatal(err.Error(), zap.Error(err))
	}

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
	if err != nil {
		l.Fatal(err.Error(), zap.Error(err))
	}

	resource.Expire(120)

	pool.MaxWait = 120 * time.Second

	dsn = fmt.Sprintf("postgres://tester:pa55word@localhost:%s/tdb?sslmode=disable", resource.GetPort("5432/tcp"))

	err = pool.Retry(func() error {
		p, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			return err
		}

		err = p.Ping(context.Background())
		if err != nil {
			return err
		}

		up, err := os.ReadFile("./testdata/setup.sql")
		if err != nil {
			return err
		}

		_, err = p.Exec(context.Background(), string(up))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		l.Fatal(err.Error(), zap.Error(err))
	}

	// Run test
	c := m.Run()

	defer func() {
		p, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			l.Fatal(err.Error(), zap.Error(err))
		}

		err = p.Ping(context.Background())
		if err != nil {
			l.Fatal(err.Error(), zap.Error(err))
		}

		down, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			l.Fatal(err.Error(), zap.Error(err))
		}

		_, err = p.Exec(context.Background(), string(down))
		if err != nil {
			l.Fatal(err.Error(), zap.Error(err))
		}

		err = pool.Purge(resource)
		if err != nil {
			l.Fatal(err.Error(), zap.Error(err))
		}
	}()

	os.Exit(c)
}

func TestNew(t *testing.T) {
	d, err := db.New(dsn)
	require.Nil(t, err)

	m := models.New(d)

	require.NotNil(t, m)
	require.NotEmpty(t, m)
}
