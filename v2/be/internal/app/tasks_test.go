package app_test

import (
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

func TestHandleCreateTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		session := scs.New()

		h := app.HandleCreateTask(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())
		body := map[string]string{"title": "running", "description": "just keeping fit"}

		ts := httptest.NewServer(session.LoadAndSave(m(h)))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.POST("/create").
			WithJSON(body).
			Expect().
			Status(http.StatusCreated).
			HasContentType("application/json")
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body map[string]string
			code int
		}{
			{
				name: "bad body",
				body: map[string]string{"name": "running"},
				code: http.StatusBadRequest,
			},
			{
				name: "invalid data",
				body: map[string]string{"title": "", "description": "just empty"},
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "duplicate data",
				body: map[string]string{"title": "test", "description": "duplicated"},
				code: http.StatusConflict,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				session := scs.New()

				h := app.HandleCreateTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, db.NewID())

				ts := httptest.NewServer(session.LoadAndSave(m(h)))
				defer ts.Close()

				e := httpexpect.Default(t, ts.URL)

				e.POST("/create").
					WithJSON(tt.body).
					Expect().
					Status(tt.code).
					HasContentType("application/json")

			})
		}
	})
}
