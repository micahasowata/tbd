package testdata

import (
	"context"
	"v2/be/internal/models"
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
