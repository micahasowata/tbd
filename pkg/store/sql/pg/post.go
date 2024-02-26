package pg

import "github.com/micahasowata/tbd/pkg/domain"

func (s *PGStore) CreatePost(post *domain.Post) (*domain.Post, error) {
	return nil, nil
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
