package domain

import (
	"time"

	"github.com/micahasowata/tbd/pkg/security"
)

type Store interface {
	CreateUser(*User) (*User, error)
	DeleteUser(int) error
	DeleteAllUsers() error
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*User, error)
}

type JWT interface {
	NewJWT(*security.Claims, time.Duration) ([]byte, error)
	VerifyJWT([]byte) (*security.Claims, error)
}
