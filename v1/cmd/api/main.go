package main

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/security"
	"github.com/micahasowata/tbd/pkg/store"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
	"github.com/pseidemann/finish"
	"go.uber.org/zap"
)

type server struct {
	*jason.Jason

	logger   *zap.Logger
	validate *validator.Validate
	store    domain.Store
	tokens   domain.JWT
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	var dsn string

	switch os.Getenv("MODE") {
	case "production":
		dsn = os.Getenv("PROD_DB_DSN")
	default:
		dsn = os.Getenv("DEV_DB_DSN")
	}

	db, err := store.New(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logger.Info("db connected successfully", zap.String("dsn", dsn))

	key := []byte(os.Getenv("KEY"))
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
		ErrorLog: zap.NewStdLog(logger),
	}

	logger.Info("server started", zap.String("port", srv.Addr))

	manager := finish.New()
	manager.Log = logger.Sugar()
	manager.Add(srv)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error(err.Error())
		}
	}()

	manager.Wait()
}
