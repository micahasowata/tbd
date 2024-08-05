package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"v2/be/internal/app"
	"v2/be/internal/app/testdata"

	"github.com/alexedwards/scs/v2"
	"github.com/gavv/httpexpect/v2"
	"go.uber.org/zap"
)

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
			JSON().Object().ContainsKey("payload")
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
					Status(tt.code)
			})
		}
	})
}
