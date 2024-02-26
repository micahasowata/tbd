package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/security"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
)

type server struct {
	*jason.Jason

	logger   *slog.Logger
	validate *validator.Validate
	store    domain.Store
	tokens   domain.JWT
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

	key := []byte(os.Getenv("TOKEN_KEY"))
	token, err := security.NewToken(key)
	if err != nil {
		panic(err)
	}

	s := &server{
		Jason:    jason.New(1_048_576, false, true),
		logger:   logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
		store:    pg.New(db),
		tokens:   token,
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
