package pg

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpPost() (*PGStore, *domain.Post, int) {
	s, u := setUpUser()
	err := s.DeleteAllUsers()
	if err != nil {
		log.Fatal(err)
		return nil, nil, 0
	}

	user, err := s.CreateUser(u)
	if err != nil {
		log.Fatal(err)
		return nil, nil, 0
	}

	post := &domain.Post{
		UserID: user.ID,
		Title:  gofakeit.BookTitle(),
		Body:   gofakeit.Paragraph(1, 10, 75, " "),
	}

	return s, post, user.ID
}
func TestCreatePost(t *testing.T) {
	s, p, _ := setUpPost()
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
	s, p, id := setUpPost()

	t.Run("valid", func(t *testing.T) {
		for range 3 {
			p, err := s.CreatePost(p)
			require.Nil(t, err)
			require.NotNil(t, p)
		}

		posts, err := s.GetUserPosts(id)
		require.Nil(t, err)
		require.NotNil(t, posts)
		require.NotEmpty(t, posts)
		assert.Equal(t, 3, len(posts))
	})

	t.Run("no posts", func(t *testing.T) {
		for range 3 {
			p, err := s.CreatePost(p)
			require.Nil(t, err)
			require.NotNil(t, p)
		}

		err := s.DeleteAllUsers()
		require.Nil(t, err)

		posts, err := s.GetUserPosts(id)
		require.Nil(t, err)
		require.Empty(t, posts)
	})
}

func TestDeletePost(t *testing.T) {

}

func TestDeleteAllPosts(t *testing.T) {

}
