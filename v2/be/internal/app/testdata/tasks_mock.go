package testdata

import (
	"context"

	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/brianvoe/gofakeit/v7"
)

type TM struct{}

func NewTM() *TM {
	return &TM{}
}

func (m *TM) Create(ctx context.Context, t *models.Task) error {
	if t.Title == "test" {
		return models.ErrDuplicateTask
	}

	if t.Title == "testX" {
		return models.ErrOpFailed
	}

	return nil
}

func (m *TM) All(ctx context.Context, userID string) ([]*models.Task, error) {
	if userID == "1" {
		return nil, nil
	}

	if userID == "25" {
		return nil, models.ErrOpFailed
	}

	t := &models.Task{
		ID:          db.NewID(),
		UserID:      userID,
		Title:       gofakeit.BookTitle(),
		Description: gofakeit.Blurb(),
		Completed:   true,
	}

	return []*models.Task{t}, nil
}

func (m *TM) GetByID(ctx context.Context, id, userID string) (*models.Task, error) {
	if id == "1" || userID == "1" {
		return nil, models.ErrRecordNotFound
	}

	if id == "25" {
		return nil, models.ErrOpFailed
	}

	if id == "345" {
		c := &models.Task{
			ID:          id,
			UserID:      userID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Blurb(),
			Completed:   true,
		}
		return c, nil
	}

	t := &models.Task{
		ID:          id,
		UserID:      userID,
		Title:       gofakeit.BookTitle(),
		Description: gofakeit.Blurb(),
	}

	return t, nil
}

func (m *TM) Update(ctx context.Context, t *models.Task) error {
	if t.Title == "test" {
		return models.ErrRecordNotFound
	}

	if t.Title == "testX" {
		return models.ErrOpFailed
	}

	return nil
}

func (m *TM) Complete(ctx context.Context, id, userID string) error {
	if id == "200" {
		return models.ErrOpFailed
	}

	return nil
}

func (m *TM) Delete(ctx context.Context, id, userID string) error {
	if id == "201" {
		return models.ErrOpFailed
	}

	return nil
}
