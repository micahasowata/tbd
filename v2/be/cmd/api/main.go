package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

const authenticatedUser = "authenticatedUser"

type application struct {
	logger   *zap.Logger
	sessions *scs.SessionManager
	models   *models.Models
}

func main() {
	logger := zap.Must(zap.NewProduction())

	dsn := os.Getenv("DSN")
	pool, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions := scs.New()
	sessions.Store = pgxstore.NewWithCleanupInterval(pool, 10*time.Hour)

	app := application{
		logger:   logger,
		sessions: sessions,
		models:   models.New(pool),
	}

	srv := &http.Server{
		Addr:           net.JoinHostPort("localhost", "4567"),
		Handler:        app.routes(),
		MaxHeaderBytes: 1_048_576,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    1 * time.Minute,
	}

	logger.Info("server start", zap.String("addr", srv.Addr))

	err = srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Error("startup error", zap.Error(err))
	}
}
