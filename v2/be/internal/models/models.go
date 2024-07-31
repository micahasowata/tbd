package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOpFailed       = errors.New("op failed")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users *UsersModel
}

func New(pool *pgxpool.Pool) *Models {
	return &Models{
		Users: &UsersModel{
			pool: pool,
		},
	}
}
