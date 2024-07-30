package pg

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/utils"
)

func setUpUser() (*PGStore, *domain.User) {
	db, _ := store.NewTestDB()

	store := New(db)

	input := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, false, false, 14),
	}

	hash, _ := utils.HashPassword(input.Password)

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	return store, user
}

func TestMain(m *testing.M) {
	s, _ := setUpUser()

	code := m.Run()

	_ = s.DeleteAllUsers()

	os.Exit(code)
}
