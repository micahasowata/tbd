package main

import (
	"log"
	"net/http"
	"os"
	"v2/be/internal/app"
	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pseidemann/finish"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		panic("dsn not set")
	}

	pool, err := db.New(dsn)
	if err != nil {
		panic(err)
	}

	m := models.New(pool)

	sessions := scs.New()
	sessions.Store = pgxstore.New(pool)

	r := app.Routes(sessions, logger, m.Users, m.Tasks)

	srv := &http.Server{
		Addr:     ":4444",
		Handler:  r,
		ErrorLog: zap.NewStdLog(logger),
	}

	logger.Info("server started", zap.String("address", "http://localhost"+srv.Addr))

	fin := finish.New()
	fin.Log = logger.Sugar()

	fin.Add(srv)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fin.Wait()
}
