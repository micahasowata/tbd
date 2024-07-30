package store

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

func NewTestDB() (*sql.DB, error) {
	dsn := "postgres://possible_bed_test:q9AfytisL1xey@localhost:4500/careful_soup_test?sslmode=disable"
	return New(dsn)
}
