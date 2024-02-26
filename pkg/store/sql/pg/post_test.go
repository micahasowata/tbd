package pg

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpPost() (*PGStore, *domain.Post) {
	s, u := setUpUser()
	err := s.DeleteAllUsers()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	user, err := s.CreateUser(u)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	post := &domain.Post{
		UserID: user.ID,
		Title:  gofakeit.BookTitle(),
		Body:   gofakeit.Paragraph(1, 10, 75, " "),
	}

	return s, post
}
func TestCreatePost(t *testing.T) {
	s, p := setUpPost()
	post, err := s.CreatePost(p)
	require.Nil(t, err)

	defer s.DeleteAllUsers()

	t.Run("valid", func(t *testing.T) {
		assert.Greater(t, post.UserID, 0)
		assert.Greater(t, post.ID, 0)
		assert.Equal(t, p.Title, post.Title)
		assert.Equal(t, p.Body, post.Body)
	})

	t.Run("invalid", func(t *testing.T) {
		s.DeleteAllUsers()

		post, err := s.CreatePost(p)
		require.NotNil(t, err)
		require.Nil(t, post)
	})
}

func TestGetUserPosts(t *testing.T) {

}

func TestDeletePost(t *testing.T) {

}

func TestDeleteAllPosts(t *testing.T) {

}
