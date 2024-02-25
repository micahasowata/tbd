package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/micahasowata/tbd/pkg/store"
)

type server struct {
	logger *slog.Logger
}

func main() {
	err := godotenv.Load(".envrc")
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := store.New(os.Getenv("DB_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logger.Info("db connected successfully")

	s := &server{
		logger: logger,
	}

	srv := &http.Server{
		Addr:     ":8000",
		Handler:  s.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("server started", slog.String("port", srv.Addr))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
}
