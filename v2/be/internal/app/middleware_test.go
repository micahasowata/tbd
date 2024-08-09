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
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const authenticatedUser = "authenticatedUser"

var userID = app.CtxKey("userID")

func setSession(t *testing.T, sessions *scs.SessionManager, ctx context.Context, val string) context.Context {
	t.Helper()

	ctx, err := sessions.Load(ctx, authenticatedUser)
	require.NoError(t, err)

	sessions.Put(ctx, authenticatedUser, val)

	_, _, err = sessions.Commit(ctx)
	require.NoError(t, err)

	return ctx
}

func lsm(t *testing.T, session *scs.SessionManager, value string) func(next http.Handler) http.Handler {
	t.Helper()

	return func(next http.Handler) http.Handler {
		return session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session.Put(r.Context(), authenticatedUser, value)

			ctx := context.WithValue(r.Context(), userID, value)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}))
	}
}

func TestRequireAuthenticatedUser(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		sessions := scs.New()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r = r.WithContext(setSession(t, sessions, r.Context(), db.NewID()))

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		m := app.RequireAuthenticatedUser(zap.NewNop(), sessions, testdata.NewUM())

		sessions.LoadAndSave(m(next)).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Equal(t, body, "OK")
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			id   string
			code int
		}{
			{
				name: "empty id",
				id:   "",
				code: http.StatusUnauthorized,
			},
			{
				name: "invalid id",
				id:   "1",
				code: http.StatusUnauthorized,
			},
			{
				name: "op failed",
				id:   "25",
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				sessions := scs.New()

				ctx := setSession(t, sessions, r.Context(), tt.id)

				r = r.WithContext(ctx)

				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(tt.id))
				})

				m := app.RequireAuthenticatedUser(zap.NewNop(), sessions, testdata.NewUM())

				sessions.LoadAndSave(m(next)).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
			})
		}
	})
}

func TestGetUserID(t *testing.T) {
	t.Parallel()

	id := db.NewID()

	ctx := context.WithValue(context.Background(), userID, id)
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r = r.WithContext(ctx)

	uid := app.GetUserID(r)

	require.Equal(t, id, uid)
}

func TestGetTaskID(t *testing.T) {
	t.Parallel()

	id := db.NewID()

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/{task_id}", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("task_id", id)

	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)

	r = r.WithContext(ctx)

	h := func(w http.ResponseWriter, r *http.Request) {
		tid := app.GetTaskID(r)
		w.Write([]byte(tid))
	}

	h(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Equal(t, id, body)
}
