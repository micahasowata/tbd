package pg

import (
	"fmt"
	"slices"
	"testing"

	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/utils"
	"github.com/rs/xid"
)

func TestCreateUser(t *testing.T) {
	// Arrange
	db, err := store.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	store := NewPGStore(db)

	input := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     "John Doe",
		Email:    fmt.Sprintf("test-%s@test.com", xid.New().String()),
		Password: "password",
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}
	// Act
	user, err = store.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if user.ID == 0 {
		t.Error("id must not be zero")
	}

	if user.Name != input.Name {
		t.Errorf("Expected: %s, Actual: %s", input.Name, user.Name)
	}

	if !slices.Equal[[]byte](user.Password, hash) {
		t.Errorf("password hash must be hashed")
	}
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	db, err := store.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	store := NewPGStore(db)

	input := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     "John Doe",
		Email:    fmt.Sprintf("test-%s@test.com", xid.New().String()),
		Password: "password",
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	user, err = store.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	//Act
	err = store.DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Expected nil, got %q", err)
	}
}
