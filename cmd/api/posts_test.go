package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	s, ts := setUpTest()
	defer s.store.DeleteAllUsers()

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

			bearer, err := setupAuth(s)
			require.Nil(t, err)

			req.Header.Set("Authorization", bearer)

			res, err := ts.Client().Do(req)
			require.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}
