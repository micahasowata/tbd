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

		r.Post("/tasks/create", app.createTask)
		r.Get("/tasks", app.allTasks)
		r.Get("/tasks/{task_id}", app.getTask)
		r.Patch("/tasks/{task_id}/update", app.updateTask)
		r.Patch("/tasks/{task_id}/complete", app.completeTask)
		r.Delete("/tasks/{task_id}", app.deleteTask)
	})
	return router
}
