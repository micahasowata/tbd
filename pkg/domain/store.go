package domain

type Store interface {
	CreateUser(user *User) (*User, error)
	DeleteUser(id int) error
}
