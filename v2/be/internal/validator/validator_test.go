package validator_test

import (
	"testing"
	"v2/be/internal/validator"

	"github.com/stretchr/testify/require"
)

func TestRequiredString(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		input := "you"

		v := validator.New()
		v.RequiredString(input, "input", validator.Required)

		require.Empty(t, v.Errors())
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "empty",
				input: "",
			},
			{
				name:  "only spaces",
				input: " ",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				v := validator.New()
				v.RequiredString(tt.input, "input", validator.Required)

				require.NotEmpty(t, v.Errors())
				require.Contains(t, v.Errors(), "input")
			})
		}
	})
}

func TestMinString(t *testing.T) {
	t.Run("bad", func(t *testing.T) {
		t.Parallel()

		input := "two"

		v := validator.New()
		v.MinString(input, 4, "input", validator.Required)

		require.NotEmpty(t, v.Errors())
		require.Contains(t, v.Errors(), "input")
	})

	t.Run("good", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "equal",
				input: "four",
			},
			{
				name:  "greater",
				input: "seven",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				v := validator.New()
				v.MinString(tt.input, 4, "input", validator.Required)

				require.Empty(t, v.Errors())
			})
		}
	})
}

func TestAddError(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		v.AddError("input", validator.Required)

		require.False(t, v.Valid())
		require.Contains(t, v.Errors(), "input")
		require.Len(t, v.Errors(), 1)
	})

	t.Run("field exists", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		v.AddError("username", validator.Required)
		v.AddError("username", "must be at least 4 characters")

		require.False(t, v.Valid())
		require.Contains(t, v.Errors(), "username")
		require.Len(t, v.Errors(), 1)
		require.Equal(t, validator.Required, v.Errors()["username"])
	})

	t.Run("nil errs", func(t *testing.T) {
		t.Parallel()

		v := &validator.Validator{}
		v.AddError("input", validator.Required)

		require.NotNil(t, v.Errors())
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		input := "E}bY^8ifURXC"

		v := validator.New()
		v.CheckPassword(input, "password")

		require.True(t, v.Valid())
	})

	t.Run("insecure", func(t *testing.T) {
		t.Parallel()

		input := "mynameis"

		v := validator.New()
		v.CheckPassword(input, "password")

		require.False(t, v.Valid())
		require.Contains(t, v.Errors(), "password")
	})
}
