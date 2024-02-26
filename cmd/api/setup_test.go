package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/security"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
)

func setUpReq(method, path, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(jason.ContentType, jason.ContentTypeJSON)
	return req, nil
}

func setupUser(s *server) (*domain.User, error) {
	err := s.store.DeleteAllUsers()
	if err != nil {
		return nil, err
	}

	u, err := s.store.CreateUser(&domain.User{
		Name:     "Joe",
		Email:    "j@doe.com",
		Password: []byte("ohh!!!ohh!!!"),
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func setupBearer(t domain.JWT, user *domain.User) (string, error) {
	token, err := t.NewJWT(&domain.Claims{
		ID:    user.ID,
		Email: user.Email,
	}, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return "Bearer " + string(token), err
}

func setUpTest() (*server, *httptest.Server) {
	db, _ := store.NewTestDB()
	token, _ := security.NewToken([]byte("=9Ha*2tME_-?xPJ_e57PEaF~UfHg6sD,"))

	srv := &server{
		Jason: jason.New(100, false, true),

		logger:   slog.Default(),
		validate: validator.New(validator.WithRequiredStructEnabled()),
		store:    pg.New(db),
		tokens:   token,
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
