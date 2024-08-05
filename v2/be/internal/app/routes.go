package app

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Routes(sessions *scs.SessionManager, logger *zap.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(sessions.LoadAndSave)

	router.Get("/", HandleHealthz(logger))

	return router
}
