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

	CreatePost(*Post) (*Post, error)
	GetUserPosts(int) ([]*Post, error)
	DeletePost(int) error
	DeleteAllPosts() error
}

type JWT interface {
	NewJWT(*Claims, time.Duration) ([]byte, error)
	VerifyJWT([]byte) (*Claims, error)
}
