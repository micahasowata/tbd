package pg

import "github.com/micahasowata/tbd/pkg/domain"

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
	return nil, nil
}

func (s *PGStore) DeletePost(id int) error {
	return nil
}

func (s *PGStore) DeleteAllPosts() error {
	return nil
}
