package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/jwt"
)

var (
	ErrInvalidKey = errors.New("key must be 32 characters long")
)

type Claims struct {
	ID    int    `json:"user_id"`
	Email string `json:"email"`
}

type Token struct {
	Key []byte
}

func NewToken(key []byte) (*Token, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	return &Token{Key: key}, nil
}

func (t *Token) NewJWT(claims *Claims, span time.Duration) ([]byte, error) {
	standardClaims := jwt.Claims{
		NotBefore: time.Now().Unix(),
		IssuedAt:  time.Now().Unix(),
		Expiry:    time.Now().Add(span).Unix(),
		Issuer:    "tbd",
		Subject:   fmt.Sprintf("user-%d-%d", claims.ID, time.Now().Unix()),
		Audience:  jwt.Audience{"tbd-ui"},
	}
	return jwt.Sign(jwt.HS256, t.Key, claims, standardClaims, jwt.MaxAge(span))
}
