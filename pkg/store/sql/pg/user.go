package pg

import "github.com/micahasowata/tbd/pkg/domain"

func (s *PGStore) CreateUser(user *domain.User) (*domain.User, error) {
	user.ID = 1
	return user, nil
}
