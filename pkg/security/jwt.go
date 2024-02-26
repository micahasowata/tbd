package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/jwt"
	"github.com/micahasowata/tbd/pkg/domain"
)

var (
	ErrInvalidKey = errors.New("key must be 32 characters long")
)

type Token struct {
	Key []byte
}

var _ domain.JWT = (*Token)(nil)

func NewToken(key []byte) (*Token, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	return &Token{Key: key}, nil
}

func (t *Token) NewJWT(claims *domain.Claims, span time.Duration) ([]byte, error) {
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

func (t *Token) VerifyJWT(token []byte) (*domain.Claims, error) {
	vt, err := jwt.Verify(jwt.HS256, t.Key, token, jwt.Leeway(2*time.Minute))
	if err != nil {
		return nil, err
	}

	c := &domain.Claims{}
	vt.Claims(c)

	return c, nil
}
