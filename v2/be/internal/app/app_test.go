package app_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"v2/be/internal/app"

	"github.com/stretchr/testify/require"
)

func readTestBody(t *testing.T, body io.Reader) string {
	t.Helper()

	b, err := io.ReadAll(body)
	require.NoError(t, err)

	return string(bytes.TrimSpace(b))
}

func TestHandleHealthz(t *testing.T) {
	t.Parallel()

	h := app.HandleHealthz()

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	h.ServeHTTP(rr, r)

	require.Equal(t, http.StatusOK, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "OK")
}
