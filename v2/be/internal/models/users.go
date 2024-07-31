package models

import (
	"context"
	"errors"
	"strings"
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
	Password []byte `json:"-"`
}

type UsersModel struct {
	pool *pgxpool.Pool
}

func (m *UsersModel) Create(ctx context.Context, u *User) error {
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

	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		switch {
		case strings.Contains(db.FormatErr(err), "users_username_key"):
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

func (m *UsersModel) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, password
	FROM users
	WHERE username = $1`

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var u User
	err = tx.QueryRow(ctx, query, username).Scan(
		&u.ID,
		&u.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *UsersModel) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return false, err
	}

	defer tx.Rollback(ctx)

	var exists bool
	err = tx.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, err
	}

	return exists, err
}
