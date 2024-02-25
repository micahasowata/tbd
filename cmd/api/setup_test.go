package main

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
)

func setUpTest() (*server, *httptest.Server) {
	db, _ := store.NewTestDB()

	srv := &server{
		Jason:    jason.New(100, false, true),
		logger:   slog.Default(),
		validate: validator.New(validator.WithRequiredStructEnabled()),
		store:    pg.New(db),
	}

	ts := httptest.NewServer(srv.routes())

	return srv, ts
}

func TestMain(m *testing.M) {
	srv, _ := setUpTest()

	code := m.Run()

	_ = srv.store.DeleteAllUsers()

	os.Exit(code)
}
