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

	return nil
}

func (m *TM) All(ctx context.Context, userID string) ([]*models.Task, error) {
	if userID == "1" {
		return nil, nil
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
