package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"v2/be/internal/app"
	"v2/be/internal/app/testdata"
	"v2/be/internal/db"

	"github.com/alexedwards/scs/v2"
	"github.com/gavv/httpexpect/v2"
	"go.uber.org/zap"
)

func lsm(session *scs.SessionManager, value string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session.Put(r.Context(), authenticatedUser, value)

			ctx := context.WithValue(r.Context(), userID, value)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}))
	}
}

func TestHandleSignup(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		logger := zap.NewNop()
		sessions := scs.New()

		body := map[string]string{"username": "alex", "password": "R#L:>t^9N?%o"}

		h := app.HandleSignup(logger, sessions, testdata.NewUM())

		ts := httptest.NewServer(sessions.LoadAndSave(h))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.POST("/signup").
			WithJSON(body).
			Expect().
			Status(http.StatusCreated).
			HasContentType("application/json").
			Cookie("session").
			HasMaxAge().Path().NotEmpty()
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body map[string]string
			code int
		}{
			{
				name: "bad body",
				body: map[string]string{"name": "Jim Carrey"},
				code: http.StatusBadRequest,
			},
			{
				name: "invalid body",
				body: map[string]string{"password": "}x0~6sz32}cN4"},
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "duplicate data",
				body: map[string]string{"username": "tester", "password": ":2~R9dq)fC9gQ"}, // Check users_mock.go for username.
				code: http.StatusConflict,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				logger := zap.NewNop()
				sessions := scs.New()

				h := app.HandleSignup(logger, sessions, testdata.NewUM())

				ts := httptest.NewServer(sessions.LoadAndSave(h))
				defer ts.Close()

				e := httpexpect.Default(t, ts.URL)

				e.POST("/signup").
					WithJSON(tt.body).
					Expect().
					Status(tt.code).
					HasContentType("application/json")
			})
		}
	})
}

func TestHandleLogin(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		logger := zap.NewNop()
		sessions := scs.New()
		body := map[string]string{"username": "alex", "password": "0~,9ZZArDp#M"}
		h := app.HandleLogin(logger, sessions, testdata.NewUM())

		ts := httptest.NewServer(sessions.LoadAndSave(h))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.POST("/login").
			WithJSON(body).
			Expect().
			Status(http.StatusOK).
			HasContentType("application/json").
			Cookie("session").HasMaxAge().Path().IsEqual("/")
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body map[string]string
			code int
		}{
			{
				name: "bad data",
				body: map[string]string{"name": "tester"},
				code: http.StatusBadRequest,
			},
			{
				name: "invalid data",
				body: map[string]string{"username": "entre", "password": "iloveyou"},
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "invalid user",
				body: map[string]string{"username": "tester", "password": "0~,9ZZArDp#M"},
				code: http.StatusNotFound,
			},
			{
				name: "invalid password",
				body: map[string]string{"username": "kolo", "password": "0~,9ZZArDp#N"},
				code: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				sessions := scs.New()

				h := app.HandleLogin(zap.NewNop(), sessions, testdata.NewUM())

				ts := httptest.NewServer(sessions.LoadAndSave(h))
				defer ts.Close()

				e := httpexpect.Default(t, ts.URL)

				e.POST("/login").
					WithJSON(tt.body).
					Expect().
					HasContentType("application/json").
					Status(tt.code)
			})
		}
	})
}

func TestHandleLogout(t *testing.T) {
	t.Parallel()

	sessions := scs.New()

	h := app.HandleLogout(zap.NewNop(), sessions)
	m := lsm(sessions, db.NewID())

	ts := httptest.NewServer(sessions.LoadAndSave(m(h)))
	defer ts.Close()

	e := httpexpect.Default(t, ts.URL)

	e.POST("/logout").
		Expect().
		HasContentType("application/json").
		Status(http.StatusOK)
}
