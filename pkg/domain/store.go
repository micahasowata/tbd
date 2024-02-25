package domain

type Store interface {
	CreateUser(*User) (*User, error)
	DeleteUser(int) error
	DeleteAllUsers() error
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*User, error)
}
