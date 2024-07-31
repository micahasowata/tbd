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
	ErrDuplicateTask = errors.New("task exist")
)

type Task struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TasksModel struct {
	pool *pgxpool.Pool
}

func (m *TasksModel) Create(ctx context.Context, t *Task) error {
	query := `INSERT INTO tasks (id, user_id, title, description)
	VALUES ($1, $2, $3, $4)`

	args := []any{t.ID, t.UserID, t.Title, t.Description}

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
		case strings.Contains(db.FormatErr(err), "tasks_title_user_id_key"):
			return ErrDuplicateTask
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
