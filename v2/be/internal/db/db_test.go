package db

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "valid dsn",
			dsn:     "postgres://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/postgres?sslmode=disable",
			wantErr: false,
		},
		{
			name:    "invalid dsn",
			dsn:     "invalid://dsn",
			wantErr: true,
		},
		{
			name:    "empty dsn",
			dsn:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := New(tt.dsn)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, pool)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, pool)

				err = pool.Ping(context.Background())
				require.NoError(t, err)
				pool.Close()
			}
		})
	}
}

func TestFormatErr(t *testing.T) {
	tests := []struct {
		name           string
		inputErr       error
		expectedOutput string
	}{
		{
			name:           "pgconn.PgError with ConstraintName",
			inputErr:       &pgconn.PgError{ConstraintName: "unique_constraint"},
			expectedOutput: "unique_constraint",
		},
		{
			name:           "non-pgconn error",
			inputErr:       errors.New("generic error"),
			expectedOutput: "",
		},
		{
			name:           "nil error",
			inputErr:       nil,
			expectedOutput: "",
		},
		{
			name:           "pgconn.PgError without ConstraintName",
			inputErr:       &pgconn.PgError{},
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatErr(tt.inputErr)
			require.Equal(t, tt.expectedOutput, result)
		})
	}
}
