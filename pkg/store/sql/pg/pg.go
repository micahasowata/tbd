package pg

import "database/sql"

type PGStore struct {
	db *sql.DB
}

func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{
		db: db,
	}
}
