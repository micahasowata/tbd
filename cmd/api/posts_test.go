package main

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	s, ts := setUpTest()
	defer func() {
		err := s.store.DeleteAllUsers()
		if err != nil {
			panic(err)
		}
	}()

	tests := []struct {
		name string
		body string
		code int
	}{
		{
			name: "valid",
			body: `{"title": "Hello World", "body": "Hello World is just a post"}`,
			code: http.StatusOK,
		},
		{
			name: "bad body",
			body: `{"title": "Hello World}`,
			code: http.StatusBadRequest,
		},
		{
			name: "bad data",
			body: `{"title": "HelloWorld"}`,
			code: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := setUpReq(http.MethodPost, ts.URL+"/v1/posts/create", tt.body)
			require.Nil(t, err)

			user, err := setupUser(s)
			require.Nil(t, err)

			bearer, err := setupBearer(s.tokens, user)
			require.Nil(t, err)

			req.Header.Set("Authorization", bearer)

			res, err := ts.Client().Do(req)
			require.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}

func TestGetPost(t *testing.T) {
	s, ts := setUpTest()
	defer func() {
		err := s.store.DeleteAllUsers()
		if err != nil {
			panic(err)
		}
	}()

	user, err := setupUser(s)
	require.Nil(t, err)

	post := &domain.Post{
		UserID: user.ID,
		Title:  "Hello",
		Body:   "Body",
	}

	post, err = s.store.CreatePost(post)
	require.Nil(t, err)

	tests := []struct {
		name string
		path string
		code int
	}{
		{
			name: "valid",
			path: ts.URL + "/v1/posts/" + strconv.Itoa(post.ID),
			code: http.StatusOK,
		},
		{
			name: "invalid id",
			path: ts.URL + "/v1/posts/" + "abc",
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid",
			path: ts.URL + "/v1/posts/" + "56",
			code: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := setUpReq(http.MethodGet, tt.path, "")
			require.Nil(t, err)

			bearer, err := setupBearer(s.tokens, user)
			require.Nil(t, err)

			req.Header.Set("Authorization", bearer)

			res, err := ts.Client().Do(req)
			require.Nil(t, err)
			require.NotNil(t, res)

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
