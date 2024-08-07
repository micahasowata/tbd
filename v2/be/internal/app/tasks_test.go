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
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
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

func TestHandleListTasks(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		session := scs.New()
		h := app.HandleListTasks(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		ts := httptest.NewServer(session.LoadAndSave(m(h)))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.GET("/all").
			Expect().
			Status(http.StatusOK).
			HasContentType("application/json").
			JSON().Object().Value("payload").Array().NotEmpty()
	})

	t.Run("errors", func(t *testing.T) {
		t.Parallel()

		session := scs.New()
		h := app.HandleListTasks(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, "1")

		ts := httptest.NewServer(session.LoadAndSave(m(h)))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.GET("/all").
			Expect().
			Status(http.StatusOK).
			HasContentType("application/json").
			JSON().Object().Value("payload").Array().IsEmpty()
	})
}

func TestHandleGetTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		tid := db.NewID()

		session := scs.New()

		h := app.HandleGetTask(zap.NewNop(), testdata.NewTM())

		m := lsm(t, session, db.NewID())

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("task_id", tid)
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

		ts := httptest.NewServer(session.LoadAndSave(m(h)))
		defer ts.Close()

		e := httpexpect.Default(t, ts.URL)

		e.GET("/tasks/{task_id}", tid).
			WithContext(ctx).
			Expect().
			Status(http.StatusFound).
			JSON().Object().Value("payload").Object().NotEmpty()
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			tid  string
			uid  string
		}{
			{
				name: "invalid task",
				tid:  "1",
				uid:  db.NewID(),
			},
			{
				name: "invalid user",
				tid:  db.NewID(),
				uid:  "1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				session := scs.New()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				h := app.HandleGetTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, tt.uid)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("task_id", tt.tid)
				ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)

				r = r.WithContext(ctx)

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)

				require.Equal(t, http.StatusNotFound, rr.Code)
			})
		}
	})
}
