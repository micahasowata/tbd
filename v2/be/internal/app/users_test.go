package app_test

import (
	"bytes"
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

func TestHandleSignup(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"username": "alex", "password": "R#L:>t^9N?%o"}`)))

		sessions := scs.New()
		h := app.HandleSignup(zap.NewNop(), sessions, testdata.NewUM())

		sessions.LoadAndSave(h).ServeHTTP(rr, r)

		require.Equal(t, http.StatusCreated, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, rs.Header.Get("Content-Type"), "application/json")
		require.Len(t, rs.Cookies(), 1)
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body string
			code int
		}{
			{
				name: "bad body",
				body: `{"name": "Jim Carrey"}`,
				code: http.StatusBadRequest,
			},
			{
				name: "invalid body",
				body: `{"password": "}x0~6sz32}cN4"}`,
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "duplicate data",
				body: `{"username": "tester", "password": ":2~R9dq)fC9gQ"}`,
				code: http.StatusConflict,
			},
			{
				name: "op failed",
				body: `{"username": "testX", "password": ":2~R9dq)fC9gQ"}`,
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))

				sessions := scs.New()

				h := app.HandleSignup(zap.NewNop(), sessions, testdata.NewUM())

				sessions.LoadAndSave(h).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
				require.Equal(t, rs.Header.Get("Content-Type"), "application/json")
				require.Len(t, rs.Cookies(), 0)
			})
		}
	})
}

func TestHandleLogin(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"username": "alex", "password": "0~,9ZZArDp#M"}`)))

		sessions := scs.New()

		h := app.HandleLogin(zap.NewNop(), sessions, testdata.NewUM())

		sessions.LoadAndSave(h).ServeHTTP(rr, r)

		require.Equal(t, http.StatusOK, rr.Code)

		rs := rr.Result()
		defer rs.Body.Close()

		body := readTestBody(t, rs.Body)

		require.Contains(t, body, "payload")
		require.Equal(t, rs.Header.Get("Content-Type"), "application/json")
		require.Len(t, rs.Cookies(), 1)
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name string
			body string
			code int
		}{
			{
				name: "bad data",
				body: `{"name": "tester"}`,
				code: http.StatusBadRequest,
			},
			{
				name: "invalid data",
				body: `{"username": "entre", "password": "iloveyou"}`,
				code: http.StatusUnprocessableEntity,
			},
			{
				name: "invalid user",
				body: `{"username": "tester", "password": "0~,9ZZArDp#M"}`,
				code: http.StatusNotFound,
			},
			{
				name: "invalid password",
				body: `{"username": "kolo", "password": "0~,9ZZArDp#N"}`,
				code: http.StatusNotFound,
			},
			{
				name: "op failed",
				body: `{"username": "testX", "password": "0~,9ZZArDp#N"}`,
				code: http.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				rr := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))

				sessions := scs.New()

				h := app.HandleLogin(zap.NewNop(), sessions, testdata.NewUM())

				sessions.LoadAndSave(h).ServeHTTP(rr, r)

				require.Equal(t, tt.code, rr.Code)

				rs := rr.Result()
				defer rs.Body.Close()

				body := readTestBody(t, rs.Body)

				require.Contains(t, body, "error")
				require.Equal(t, rs.Header.Get("Content-Type"), "application/json")
				require.Len(t, rs.Cookies(), 0)
			})
		}
	})
}

func TestHandleLogout(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	sessions := scs.New()

	h := app.HandleLogout(zap.NewNop(), sessions)
	m := lsm(t, sessions, db.NewID())

	sessions.LoadAndSave(m(h)).ServeHTTP(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "payload")
	require.Equal(t, rs.Header.Get("Content-Type"), "application/json")
	require.Len(t, rs.Cookies(), 1)
}
