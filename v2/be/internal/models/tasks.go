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

func (m *TasksModel) All(ctx context.Context, userID string) ([]*Task, error) {
	query := `SELECT id, title, description, completed
	FROM tasks
	WHERE user_id = $1`

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var tasks []*Task

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task
		terr := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Completed,
		)

		if terr != nil {
			return nil, terr
		}

		tasks = append(tasks, &t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *TasksModel) GetByID(ctx context.Context, id, userID string) (*Task, error) {
	query := `SELECT id, title, description, completed
	FROM tasks
	WHERE id = $1 AND user_id = $2`

	args := []any{id, userID}

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var t Task

	err = tx.QueryRow(ctx, query, args...).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Completed,
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

	return &t, nil
}

func (m *TasksModel) Update(ctx context.Context, t *Task) error {
	query := `UPDATE tasks
	SET title = $1, description = $2
	WHERE id = $3 AND completed = false`

	args := []any{t.Title, t.Description, t.ID}

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
