package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/micahasowata/tbd/pkg/store"
)

func main() {
	err := godotenv.Load(".envrc")
	if err != nil {
		panic(err)
	}

	db, err := store.New(os.Getenv("DB_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	slog.Info("Hello world!")
}
