package parser_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"v2/be/internal/parser"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		w := httptest.NewRecorder()

		err := parser.Write(w, http.StatusOK, parser.Envelope{"testing": "yes"})

		require.Nil(t, err)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad", func(t *testing.T) {
		w := httptest.NewRecorder()

		err := parser.Write(w, http.StatusOK, parser.Envelope{"age": complex128(2)})

		require.NotNil(t, err)
	})
}
