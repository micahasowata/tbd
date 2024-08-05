package testdata

import (
	"context"
	"errors"
	"v2/be/internal/models"
)

// UM is short for Users Mock
type UM struct {
}

func NewUM() *UM {
	return &UM{}
}

func (m *UM) Create(ctx context.Context, u *models.User) error {
	switch u.Username {
	case "":
		return errors.New("no username")
	case "tester":
		return models.ErrDuplicateUsername
	}

	return nil
}
