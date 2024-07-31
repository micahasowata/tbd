package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(app.sessions.LoadAndSave)

	router.Post("/signup", app.signup)
	router.Post("/login", app.login)
	return router
}
