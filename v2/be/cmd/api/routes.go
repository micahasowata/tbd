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

	router.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)

		r.Post("/logout", app.logout)
	})
	return router
}
