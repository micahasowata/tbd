package models

import (
	"context"
	"errors"
	"v2/be/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateUsername = errors.New("username exists")
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type UsersModel struct {
	pool *pgxpool.Pool
}

func (m *UsersModel) Create(ctx context.Context, u User) error {
	query := `INSERT INTO users (id, username, password)
	VALUES ($1, $2, $3)`

	args := []any{u.ID, u.Username, u.Password}

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		switch {
		case db.FormatErr(err) == "users_username_key":
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	if result.RowsAffected() != 1 {
		return ErrOpFailed
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
