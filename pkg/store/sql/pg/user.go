package pg

import (
	"github.com/micahasowata/tbd/pkg/domain"
)

func (s *PGStore) CreateUser(user *domain.User) (*domain.User, error) {
	query := `
	INSERT INTO users (name, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, name, email, password`

	args := []any{user.Name, user.Email, user.Password}

	err := s.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
