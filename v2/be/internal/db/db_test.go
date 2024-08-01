package db_test

import (
	"errors"
	"testing"
	"v2/be/internal/db"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d, err := db.New("postgresql://user:secret@localhost:5432/db?sslmode=disable")
		require.Nil(t, err)

		require.NotNil(t, d)
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
