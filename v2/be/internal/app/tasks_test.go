package app_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"v2/be/internal/app"
	"v2/be/internal/app/testdata"
	"v2/be/internal/db"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setTaskID(t *testing.T, id string) context.Context {
	t.Helper()

	rtx := chi.NewRouteContext()
	rtx.URLParams.Add("task_id", id)
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rtx)

	return ctx
}

func TestHandleCreateTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"title": "running", "description": "just keeping fit"}`)))

		session := scs.New()

		h := app.HandleCreateTask(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusCreated, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body string
			code int
		}{
			{
				name: "bad body",
				body: `{"name": "running"}`,
				code: http.StatusBadRequest,
			},
			{
				name: "invalid data",
				body: `{"title": "", "description": "just empty"}`,
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "duplicate data",
				body: `{"title": "test", "description": "duplicated"}`,
				code: http.StatusConflict,
			},
			{
				name: "op failed",
				body: `{"title": "testX", "description": "duplicated"}`,
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))

				session := scs.New()

				h := app.HandleCreateTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, db.NewID())

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
				require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
			})
		}
	})
}

func TestHandleListTasks(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		session := scs.New()
		h := app.HandleListTasks(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		session := scs.New()
		h := app.HandleListTasks(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, "1")

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "[]")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		session := scs.New()
		h := app.HandleListTasks(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, "25")

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusInternalServerError, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "error")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})
}

func TestHandleGetTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := setTaskID(t, db.NewID())
		r = r.WithContext(ctx)

		session := scs.New()

		h := app.HandleGetTask(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)
		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			tid  string
			code int
		}{
			{
				name: "invalid task",
				tid:  "1",
				code: http.StatusNotFound,
			},
			{
				name: "op failed",
				tid:  "25",
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				session := scs.New()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx := setTaskID(t, tt.tid)
				r = r.WithContext(ctx)

				h := app.HandleGetTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, db.NewID())

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)
				require.Contains(t, body, "error")
				require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
			})
		}
	})
}

func TestHandleUpdateTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(`{"title": "read", "description": "complete a chapter"}`)))
		ctx := setTaskID(t, db.NewID())
		r = r.WithContext(ctx)

		session := scs.New()

		h := app.HandleUpdateTask(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)
		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			tid  string
			body string
			code int
		}{
			{
				name: "bad body",
				tid:  db.NewID(),
				body: `{"name":"happier"}`,
				code: http.StatusBadRequest,
			},
			{
				name: "invalid data",
				tid:  db.NewID(),
				body: `{"title":"", "description":"early morning run"}`,
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "missing data",
				tid:  "1",
				body: `{"title":"learn testing", "description":"practice TDD"}`,
				code: http.StatusNotFound,
			},
			{
				name: "op failed",
				tid:  "25",
				body: `{"title":"learn testing", "description":"practice TDD"}`,
				code: http.StatusInternalServerError,
			},
			{
				name: "completed task",
				tid:  "345",
				body: `{"title":"learn testing", "description":"practice TDD"}`,
				code: http.StatusNotModified,
			},
			{
				name: "task not found",
				tid:  db.NewID(),
				body: `{"title":"test", "description":"update fails"}`,
				code: http.StatusNotFound,
			},
			{
				name: "op failed",
				tid:  db.NewID(),
				body: `{"title":"testX", "description":"update fails"}`,
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(tt.body)))
				ctx := setTaskID(t, tt.tid)
				r = r.WithContext(ctx)

				session := scs.New()

				h := app.HandleUpdateTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, db.NewID())

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
				require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
			})
		}
	})
}

func TestHandleCompleteTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodPatch, "/", nil)
		ctx := setTaskID(t, db.NewID())
		r = r.WithContext(ctx)

		h := app.HandleCompleteTask(zap.NewNop(), testdata.NewTM())

		session := scs.New()
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)
		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			tid  string
			code int
		}{
			{
				name: "missing task",
				tid:  "1",
				code: http.StatusNotFound,
			},
			{
				name: "task not found",
				tid:  "25",
				code: http.StatusInternalServerError,
			},
			{
				name: "completed task",
				tid:  "345",
				code: http.StatusNotModified,
			},
			{
				name: "completion error",
				tid:  "200",
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()

				r := httptest.NewRequest(http.MethodPatch, "/", nil)
				ctx := setTaskID(t, tt.tid)
				r = r.WithContext(ctx)

				h := app.HandleCompleteTask(zap.NewNop(), testdata.NewTM())

				session := scs.New()
				m := lsm(t, session, db.NewID())

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)
				require.Equal(t, tt.code, rr.Code)
			})
		}
	})
}

func TestHandleDeleteTask(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		ctx := setTaskID(t, db.NewID())
		r = r.WithContext(ctx)

		session := scs.New()

		h := app.HandleDeleteTask(zap.NewNop(), testdata.NewTM())
		m := lsm(t, session, db.NewID())

		session.LoadAndSave(m(h)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			tid  string
			code int
		}{
			{
				name: "missing task",
				tid:  "1",
				code: http.StatusNotFound,
			},
			{
				name: "get error",
				tid:  "25",
				code: http.StatusInternalServerError,
			},
			{
				name: "delete error",
				tid:  "201",
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()

				r := httptest.NewRequest(http.MethodDelete, "/", nil)
				ctx := setTaskID(t, tt.tid)
				r = r.WithContext(ctx)

				session := scs.New()
				h := app.HandleDeleteTask(zap.NewNop(), testdata.NewTM())
				m := lsm(t, session, db.NewID())

				session.LoadAndSave(m(h)).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
				require.Equal(t, "application/json", rs.Header.Get("Content-Type"))
			})
		}
	})
}
