package parser_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"v2/be/internal/parser"

	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	input := "t<>yet"
	want := "t&lt;&gt;yet"
	got := parser.Sanitize(input)

	require.Equal(t, want, got)
}

func TestRead(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d := `{"name": "input"}`
		var i struct {
			Name string `json:"name"`
		}

		w := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(d))
		require.Nil(t, err)

		err = parser.Read(w, r, &i)
		require.Nil(t, err)

		require.Equal(t, "input", i.Name)
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name string
			body string
		}{
			{
				name: "syntax error",
				body: `{"name": "input",}`,
			},
			{
				name: "badly formed body",
				body: `{"name": "input"`,
			},
			{
				name: "type error",
				body: `{"name": 234}`,
			},
			{
				name: "no field type error",
				body: `["foo", "bar"]`,
			},
			{
				name: "empty",
				body: ``,
			},
			{
				name: "unknown",
				body: `{"age": 234}`,
			},
			{
				name: "double bodies",
				body: `{"name":"go"}{"name":"input"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()

				r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))
				require.Nil(t, err)

				var input struct {
					Name string `json:"name"`
				}

				err = parser.Read(w, r, &input)
				require.NotNil(t, err)
			})
		}
	})
}
