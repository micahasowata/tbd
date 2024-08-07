package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"v2/be/internal/app"
	"v2/be/internal/app/testdata"
	"v2/be/internal/db"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const authenticatedUser = "authenticatedUser"

func TestRequireAuthenticatedUser(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		sessions := scs.New()

		ctx, err := sessions.Load(r.Context(), authenticatedUser)
		require.NoError(t, err)

		sessions.Put(ctx, authenticatedUser, db.NewID())

		_, _, err = sessions.Commit(ctx)
		require.NoError(t, err)

		r = r.WithContext(ctx)

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
		}{
			{
				name: "empty id",
				id:   "",
			},
			{
				name: "invalid id",
				id:   "1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				sessions := scs.New()

				ctx, err := sessions.Load(r.Context(), authenticatedUser)
				require.NoError(t, err)

				sessions.Put(ctx, authenticatedUser, tt.id)

				_, _, err = sessions.Commit(ctx)
				require.NoError(t, err)

				r = r.WithContext(ctx)

				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(tt.id))
				})

				m := app.RequireAuthenticatedUser(zap.NewNop(), sessions, testdata.NewUM())

				sessions.LoadAndSave(m(next)).ServeHTTP(rr, r)

				require.Equal(t, http.StatusUnauthorized, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
			})
		}
	})
}
