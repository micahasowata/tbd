package security

import (
	"errors"
)

var (
	ErrInvalidKey = errors.New("key must be 32 characters long")
)

type Claims struct {
	ID    int
	Email string
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
