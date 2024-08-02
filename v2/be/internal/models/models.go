package models

import (
	"errors"
	"v2/be/internal/db"
)

var (
	ErrOpFailed       = errors.New("op failed")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users *UsersModel
	Tasks *TasksModel
}

func New(db db.DB) *Models {
	return &Models{
		Users: &UsersModel{
			pool: db,
		},
		Tasks: &TasksModel{
			pool: db,
		},
	}
}
