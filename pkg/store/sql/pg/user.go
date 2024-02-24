package pg

import (
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/utils"
)

func (s *PGStore) CreateUser(user *domain.User) (*domain.User, error) {
	query := `
	INSERT INTO users (name, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, name, email, password`

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	args := []any{user.Name, user.Email, hash}

	err = s.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&hash,
	)

	if err != nil {
		return nil, err
	}

	user.Password = string(hash)

	return user, nil
}
