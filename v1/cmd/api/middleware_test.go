package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthUser(t *testing.T) {
	s, _ := setUpTest()
	defer func() {
		err := s.store.DeleteAllUsers()
		if err != nil {
			panic(err)
		}
	}()
	u, err := s.store.CreateUser(&domain.User{
		Name:     "John",
		Email:    "j@doe.co",
		Password: []byte("oops!!oo!!@"),
	})

	require.Nil(t, err)

	claims := &domain.Claims{
		ID:    u.ID,
		Email: u.Email,
	}

	token, err := s.tokens.NewJWT(claims, 3*time.Hour)
	require.Nil(t, err)

	tests := []struct {
		name   string
		token  []byte
		header string
		code   int
	}{
		{
			name:   "valid token",
			token:  token,
			header: "Bearer " + string(token),
			code:   http.StatusOK,
		},
		{
			name:   "no token",
			token:  []byte{},
			header: "Bearer",
			code:   http.StatusUnauthorized,
		},
		{
			name:   "no token",
			token:  token,
			header: "Basic " + string(token),
			code:   http.StatusUnauthorized,
		},
		{
			name:   "invalid token",
			token:  slices.Concat[[]byte](token, []byte("invalid")),
			header: "Bearer " + string(slices.Concat[[]byte](token, []byte("invalid"))),
			code:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodPost, "/", nil)
			require.Nil(t, err)

			r.Header.Set("Authorization", tt.header)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("OK"))
				require.Nil(t, err)
			})

			s.authUser(next).ServeHTTP(rr, r)

			rs := rr.Result()

			assert.Equal(t, tt.code, rs.StatusCode)
		})
	}
}

func TestAuthUser_InvalidUser(t *testing.T) {
	s, _ := setUpTest()

	u, err := s.store.CreateUser(&domain.User{
		Name:     "John",
		Email:    "j@doe.co",
		Password: []byte("oops!!oo!!@"),
	})

	require.Nil(t, err)

	claims := &domain.Claims{
		ID:    u.ID,
		Email: u.Email,
	}

	token, err := s.tokens.NewJWT(claims, 3*time.Hour)
	require.Nil(t, err)

	err = s.store.DeleteUser(u.ID)
	require.Nil(t, err)

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodPost, "/", nil)
	require.Nil(t, err)

	r.Header.Set("Authorization", "Bearer "+string(token))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		require.Nil(t, err)
	})

	s.authUser(next).ServeHTTP(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusForbidden, rs.StatusCode)
}
