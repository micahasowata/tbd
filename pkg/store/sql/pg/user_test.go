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
	createdUser, err := store.CreateUser(user)
	if err != nil {
		t.Fatal()
	}

	// Assert
	if createdUser.ID == 0 {
		t.Error("id must not be zero")
	}

	if createdUser.Name != user.Name {
		t.Errorf("Expected: %s, Actual: %s", user.Name, createdUser.Name)
	}

	if createdUser.Password == user.Password {
		t.Error("password must be hashed")
	}
}
