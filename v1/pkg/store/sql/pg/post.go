package pg

import (
	"database/sql"
	"errors"

	"github.com/micahasowata/tbd/pkg/domain"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

func (s *PGStore) CreatePost(post *domain.Post) (*domain.Post, error) {
	query := `
	INSERT INTO posts (user_id,title, body)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, user_id, title, body`

	args := []any{post.UserID, post.Title, post.Body}

	err := s.db.QueryRow(query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UserID,
		&post.Title,
		&post.Body,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}
func (s *PGStore) GetUserPosts(userID int) ([]*domain.Post, error) {
	query := `
	SELECT id, created_at, user_id, title, body
	FROM posts
	WHERE user_id = $1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	posts := []*domain.Post{}

	for rows.Next() {
		var post domain.Post
		err = rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.UserID,
			&post.Title,
			&post.Body,
		)

		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PGStore) DeletePost(post *domain.Post) error {
	query := `
	DELETE FROM posts
	WHERE id = $1
	AND user_id = $2`

	args := []any{post.ID, post.UserID}

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *PGStore) GetPost(post *domain.Post) (*domain.Post, error) {
	query := `
	SELECT id, created_at, user_id, title, body
	FROM posts
	WHERE id = $1
	AND user_id = $2`

	args := []any{post.ID, post.UserID}

	err := s.db.QueryRow(query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UserID,
		&post.Title,
		&post.Body,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrPostNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}
