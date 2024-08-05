package app_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"v2/be/internal/app"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServerError(t *testing.T) {
	w := httptest.NewRecorder()
	logger := zap.NewNop()
	app.ServerError(w, logger, errors.New("random error"))

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
