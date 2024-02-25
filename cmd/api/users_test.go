package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpTest(t *testing.T) *httptest.Server {
	db, err := store.NewTestDB()
	require.Nil(t, err)

	srv := &server{
		logger: slog.Default(),
		store:  pg.New(db),
	}

	ts := httptest.NewServer(srv.routes())

	return ts
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		body string
		code int
	}{
		{
			name: "valid",
			body: `{"name": "tbd", "email":"me@tbd.com", "password":"://me@tbd.com://"}`,
			code: http.StatusOK,
		},
		{
			name: "bad request body",
			body: `{"name": "tbd`,
			code: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setUpTest(t)
			defer ts.Close()

			res, err := ts.Client().Post(ts.URL+"/v1/users/create", jason.ContentTypeJSON, strings.NewReader(tt.body))
			require.Nil(t, err)

			assert.Equal(t, tt.code, res.StatusCode)
		})
	}
}