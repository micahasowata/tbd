package app

import (
	"net/http"
	"v2/be/internal/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Routes(sessions *scs.SessionManager, logger *zap.Logger, u *models.UsersModel) http.Handler {
	router := chi.NewRouter()
	router.Use(sessions.LoadAndSave)

	router.Get("/", HandleHealthz())
	router.Post("/signup", HandleSignup(logger, sessions, u))
	return router
}
