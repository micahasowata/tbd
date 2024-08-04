package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func FormatErr(err error) string {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.ConstraintName
	}

	return ""
}

func NewID() string {
	u := uuid.Must(uuid.NewV7())
	return u.String()
}
