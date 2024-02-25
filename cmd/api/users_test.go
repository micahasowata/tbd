package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpTest(t *testing.T) (*server, *httptest.Server) {
	db, err := store.NewTestDB()
	require.Nil(t, err)

	srv := &server{
		logger: slog.Default(),
		store:  pg.New(db),
	}

	ts := httptest.NewServer(srv.routes())

	return srv, ts
}

func TestCreateUser(t *testing.T) {
	body := fmt.Sprintf(`{"name":%s, "email":%s, "password":%s}`,
		gofakeit.Name(), gofakeit.Email(), gofakeit.Password(true, true, true, false, false, 14))

	srv, ts := setUpTest(t)
	defer ts.Close()

	res, err := ts.Client().Post(ts.URL+"/v1/users/create", jason.ContentTypeJSON, strings.NewReader(body))
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.HTTPBodyContains(t, srv.createUser, http.MethodPost, "/v1/users/create", nil, "OK")
}
