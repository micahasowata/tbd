package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
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
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	db, err := store.New("postgres://main:HmgJYuBHO23pGp7YrHaY@tbd.cjmoua262vpy.eu-north-1.rds.amazonaws.com:5432/tbd")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logger.Info("db connected successfully")

	key := []byte("Y_1,5a?gP^M5*k3#xxjs7muWJEyGm>su")
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
