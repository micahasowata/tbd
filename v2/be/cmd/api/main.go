package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type application struct {
	logger *zap.Logger
}

func main() {
	logger := zap.Must(zap.NewProduction())

	app := application{
		logger: logger,
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

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Error("startup error", zap.Error(err))
	}
}
