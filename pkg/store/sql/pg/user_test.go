package pg

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpUser(t *testing.T) (*PGStore, *domain.User) {
	t.Helper()

	db, err := store.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

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

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	return store, user
}

func TestCreateUser(t *testing.T) {
	s, u := setUpUser(t)
	defer s.DeleteAllUsers()
	t.Run("valid", func(t *testing.T) {
		user, err := s.CreateUser(u)
		require.Nil(t, err)

		assert.Equal(t, u.Name, user.Name)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Password, user.Password)
	})

	t.Run("invalid", func(t *testing.T) {
		user, err := s.CreateUser(u)
		require.NotNil(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, ErrEmailInUse.Error())
	})

}

func TestDeleteUser(t *testing.T) {
	s, u := setUpUser(t)
	defer s.DeleteAllUsers()
	t.Run("valid", func(t *testing.T) {
		user, err := s.CreateUser(u)
		require.Nil(t, err)

		err = s.DeleteUser(user.ID)
		require.Nil(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		err := s.DeleteUser(u.ID)
		require.NotNil(t, err)

		assert.EqualError(t, err, ErrUserNotFound.Error())
	})

}

func TestGetUserByEmail(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		s, u := setUpUser(t)
		defer s.DeleteAllUsers()

		user, err := s.CreateUser(u)
		require.Nil(t, err)

		u, err = s.GetUserByEmail(user.Email)
		require.Nil(t, err)

		assert.Equal(t, u.ID, user.ID)
		assert.Equal(t, u.Name, user.Name)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Password, user.Password)
	})

	t.Run("invalid", func(t *testing.T) {
		s, u := setUpUser(t)

		user, err := s.GetUserByEmail(u.Email)
		require.NotNil(t, err)
		require.Nil(t, user)

	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		s, u := setUpUser(t)
		user, err := s.CreateUser(u)
		defer s.DeleteAllUsers()

		require.Nil(t, err)

		u, err = s.GetUserByID(user.ID)
		require.Nil(t, err)

		assert.Equal(t, u.ID, user.ID)
		assert.Equal(t, u.Name, user.Name)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Password, user.Password)
	})

	t.Run("invalid", func(t *testing.T) {
		s, u := setUpUser(t)

		user, err := s.GetUserByID(u.ID)
		require.NotNil(t, err)
		require.Nil(t, user)

	})
}

func TestDeleteAllUsers(t *testing.T) {
	s, u := setUpUser(t)
	user, err := s.CreateUser(u)
	require.Nil(t, err)

	err = s.DeleteAllUsers()
	require.Nil(t, err)

	user, err = s.GetUserByEmail(user.Email)
	require.NotNil(t, err)
	require.Nil(t, user)
}
