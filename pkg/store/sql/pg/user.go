package pg

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/micahasowata/tbd/pkg/domain"
)

var (
	DuplicateData   = "pq: duplicate key value violates unique constraint"
	ErrEmailInUse   = errors.New("email is already in use")
	ErrUserNotFound = errors.New("user not found")
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
		switch {
		case strings.Contains(err.Error(), DuplicateData):
			return nil, ErrEmailInUse
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *PGStore) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password
	FROM users
	WHERE email = $1`

	user := &domain.User{}

	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *PGStore) GetUserByID(id int) (*domain.User, error) {
	query := `SELECT id, name, email, password
	FROM users
	WHERE id = $1`

	user := &domain.User{}

	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *PGStore) DeleteUser(id int) error {
	query := `
	DELETE FROM users
	WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return ErrUserNotFound
	}

	return nil
}
