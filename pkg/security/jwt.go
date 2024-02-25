package security

import (
	"errors"
	"time"

	"github.com/kataras/jwt"
)

var (
	ErrInvalidKey = errors.New("key must be 32 characters long")
)

type Claims struct {
	ID       int    `json:"user_id"`
	Email    string `json:"email"`
	IssuedAt int64  `json:"issued_at"`
	Expired  int64  `json:"expired"`
	Issuer   string `json:"issuer"`
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

func (t *Token) NewJWT(c *Claims, span time.Duration) ([]byte, error) {
	c.IssuedAt = time.Now().Unix()
	c.Expired = time.Now().Add(span).Unix()
	c.Issuer = "tbd"

	return jwt.Sign(jwt.HS256, t.Key, c, jwt.MaxAge(span))
}
