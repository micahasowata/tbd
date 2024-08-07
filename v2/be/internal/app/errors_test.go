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
	t.Parallel()

	rr := httptest.NewRecorder()
	logger := zap.NewNop()

	app.ServerError(rr, logger, errors.New("server error"))

	require.Equal(t, http.StatusInternalServerError, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "error")
}

func TestReadError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	logger := zap.NewNop()

	app.ReadError(rr, logger, errors.New("read error"))

	require.Equal(t, http.StatusBadRequest, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "read error")
}

func TestInvalidDataError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()

	app.InvalidDataError(rr, map[string]string{"input": "should not be empty"})
	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "input")
}

func TestDuplicateDataError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	logger := zap.NewNop()

	app.DuplicateDataError(rr, logger, errors.New("duplicate data"))
	require.Equal(t, http.StatusConflict, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "duplicate data")
}

func TestMissingDataError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	logger := zap.NewNop()

	app.MissingDataError(rr, logger, errors.New("missing data"))
	require.Equal(t, http.StatusNotFound, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "missing data")
}

func TestUnauthorizedError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()

	app.UnauthorizedAccessError(rr, zap.NewNop(), errors.New("unauthorized access"))
	require.Equal(t, http.StatusUnauthorized, rr.Code)

	rs := rr.Result()
	defer rs.Body.Close()

	body := readTestBody(t, rs.Body)

	require.Contains(t, body, "unauthorized access")
}
