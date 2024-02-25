package pg

import (
	"database/sql"

	"github.com/micahasowata/tbd/pkg/domain"
)

type PGStore struct {
	db *sql.DB
}

var _ domain.Store = (*PGStore)(nil)

func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{
		db: db,
	}
}
