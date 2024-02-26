package security

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validKey   = []byte("jJn5R79QyT4vep6c2sXxG8UaKbLwN3zH")
	invalidKey = []byte("short")
)

func TestNewToken(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		token, err := NewToken(validKey)
		require.Nil(t, err)
		require.NotNil(t, token)

		assert.Equal(t, token.Key, validKey)
	})

	t.Run("invalid", func(t *testing.T) {
		token, err := NewToken(invalidKey)
		require.NotNil(t, err)
		require.Nil(t, token)

		assert.EqualError(t, err, ErrInvalidKey.Error())
	})
}

func TestNewJWT(t *testing.T) {
	token, err := NewToken(validKey)
	require.Nil(t, err)

	claims := &Claims{
		ID:    1,
		Email: "tbd@tbd",
	}

	jwt, err := token.NewJWT(claims, 15*time.Minute)
	require.Nil(t, err)
	require.NotNil(t, jwt)
}

func TestVerifyToken(t *testing.T) {
	token, err := NewToken(validKey)
	require.Nil(t, err)

	claims := &Claims{
		ID:    1,
		Email: "tbd@tbd",
	}

	jwt, err := token.NewJWT(claims, 15*time.Minute)
	require.Nil(t, err)
	require.NotNil(t, jwt)

	t.Run("valid", func(t *testing.T) {
		verifiedClaims, err := token.VerifyJWT(jwt)
		require.Nil(t, err)
		require.NotNil(t, claims)

		assert.Equal(t, verifiedClaims.ID, claims.ID)
		assert.Equal(t, verifiedClaims.Email, claims.Email)
	})

	t.Run("invalid", func(t *testing.T) {
		invalidClaims, err := token.VerifyJWT(slices.Concat[[]byte](jwt, []byte("invalid")))
		require.NotNil(t, err)
		require.Nil(t, invalidClaims)
	})
}
