package domain

import (
	"time"
)

type Store interface {
	CreateUser(*User) (*User, error)
	DeleteUser(int) error
	DeleteAllUsers() error
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*User, error)
}

type JWT interface {
	NewJWT(*Claims, time.Duration) ([]byte, error)
	VerifyJWT([]byte) (*Claims, error)
}
