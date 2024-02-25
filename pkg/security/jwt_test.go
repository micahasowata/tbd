package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		token, err := NewToken([]byte("jJn5R79QyT4vep6c2sXxG8UaKbLwN3zH"))
		require.Nil(t, err)
		require.NotNil(t, token)

		assert.Equal(t, token.Key, []byte("jJn5R79QyT4vep6c2sXxG8UaKbLwN3zH"))
	})

	t.Run("invalid", func(t *testing.T) {
		token, err := NewToken([]byte("short"))
		require.NotNil(t, err)
		require.Nil(t, token)

		assert.EqualError(t, err, ErrInvalidKey.Error())
	})
}
