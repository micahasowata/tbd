package testdata

import (
	"context"
	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/alexedwards/argon2id"
)

// UM is short for Users Mock
type UM struct {
}

func NewUM() *UM {
	return &UM{}
}

func (m *UM) Create(ctx context.Context, u *models.User) error {
	if u.Username == "tester" {
		return models.ErrDuplicateUsername
	}

	return nil
}

func (m *UM) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "tester" {
		return nil, models.ErrRecordNotFound
	}

	hash, err := argon2id.CreateHash("0~,9ZZArDp#M", argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:       db.NewID(),
		Username: username,
		Password: []byte(hash),
	}, nil
}
