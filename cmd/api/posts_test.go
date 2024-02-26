package main

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	s, ts := setUpTest()

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

	defer s.store.DeleteAllUsers()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/posts/create", strings.NewReader(tt.body))
			require.Nil(t, err)

			err = s.store.DeleteAllUsers()
			require.Nil(t, err)

			u, err := s.store.CreateUser(&domain.User{
				Name:     "Joe",
				Email:    "j@doe.com",
				Password: []byte("ohh!!!ohh!!!"),
			})
			require.Nil(t, err)

			token, err := s.tokens.NewJWT(&domain.Claims{
				ID:    u.ID,
				Email: u.Email,
			}, 15*time.Minute)
			require.Nil(t, err)

			req.Header.Set(jason.ContentType, jason.ContentTypeJSON)
			req.Header.Set("Authorization", "Bearer "+string(token))

			res, err := ts.Client().Do(req)
			require.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
