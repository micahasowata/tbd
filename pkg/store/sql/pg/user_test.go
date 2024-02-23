package pg

import (
	"testing"

	"github.com/micahasowata/tbd/pkg/domain"
)

func TestCreateUser(t *testing.T) {
	// Arrange
	store := NewPGStore(nil)
	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "John Doe",
	}

	// Act
	user, err := store.CreateUser(user)
	if err != nil {
		t.Fatal()
	}

	// Assert
	if user.ID == 0 {
		t.Errorf("id must not be zero")
	}
}
