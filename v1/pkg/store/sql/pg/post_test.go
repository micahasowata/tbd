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

	defer func() {
		err := s.DeleteAllUsers()
		require.Nil(t, err)

	}()

	t.Run("valid", func(t *testing.T) {
		assert.Greater(t, post.UserID, 0)
		assert.Greater(t, post.ID, 0)
		assert.Equal(t, p.Title, post.Title)
		assert.Equal(t, p.Body, post.Body)
	})

	t.Run("invalid", func(t *testing.T) {
		err := s.DeleteAllUsers()
		require.Nil(t, err)

		post, err := s.CreatePost(p)
		require.NotNil(t, err)
		require.Nil(t, post)
	})
}

func TestGetUserPosts(t *testing.T) {
	s, p := setUpPost()

	t.Run("valid", func(t *testing.T) {
		for range 3 {
			p, err := s.CreatePost(p)
			require.Nil(t, err)
			require.NotNil(t, p)
		}

		posts, err := s.GetUserPosts(p.UserID)
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

		posts, err := s.GetUserPosts(p.UserID)
		require.Nil(t, err)
		require.Empty(t, posts)
	})
}

func TestDeletePost(t *testing.T) {
	s, p := setUpPost()

	p, err := s.CreatePost(p)
	require.Nil(t, err)
	require.NotNil(t, p)

	err = s.DeletePost(p)
	require.Nil(t, err)

	p, err = s.GetPost(p)
	require.NotNil(t, err)
	require.Nil(t, p)
}

func TestGetPost(t *testing.T) {
	s, p := setUpPost()

	t.Run("valid", func(t *testing.T) {
		p, err := s.CreatePost(p)
		require.Nil(t, err)
		require.NotNil(t, p)

		cp, err := s.GetPost(p)
		require.Nil(t, err)
		require.NotNil(t, cp)

		assert.Equal(t, p.ID, cp.ID)
		assert.Equal(t, p.UserID, cp.UserID)
		assert.Equal(t, p.CreatedAt, cp.CreatedAt)
		assert.Equal(t, p.Title, cp.Title)
		assert.Equal(t, p.Body, cp.Body)
	})

	t.Run("invalid", func(t *testing.T) {
		p, err := s.CreatePost(p)
		require.Nil(t, err)
		require.NotNil(t, p)

		err = s.DeletePost(p)
		require.Nil(t, err)

		cp, err := s.GetPost(p)
		require.NotNil(t, err)
		require.Nil(t, cp)
	})
}
