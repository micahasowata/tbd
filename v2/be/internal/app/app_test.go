package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"v2/be/internal/app"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandleHealthz(t *testing.T) {
	logger := zap.NewNop()

	hf := app.HandleHealthz(logger)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	hf.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
}
