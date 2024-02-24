package domain

type Store interface {
	CreateUser(*User) (*User, error)
	DeleteUser(int) error
}
