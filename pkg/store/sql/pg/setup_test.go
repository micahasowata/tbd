package pg

import (
	"database/sql"
	"os"
	"testing"

	"github.com/micahasowata/tbd/pkg/store"
)

var (
	tdb *sql.DB
)

func TestMain(m *testing.M) {
	dsn := "postgres://possible_bed_test:q9AfytisL1xey@localhost:4500/careful_soup_test?sslmode=disable"
	tdb, err := store.New(dsn)
	if err != nil {
		panic(err)
	}
	defer tdb.Close()

	os.Exit(m.Run())
}
