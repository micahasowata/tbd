package app

import (
	"net/http"

	"v2/be/internal/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Routes(
	sessions *scs.SessionManager,
	logger *zap.Logger,
	u *models.UsersModel,
	t *models.TasksModel,
) http.Handler {
	router := chi.NewRouter()
	router.Use(sessions.LoadAndSave)

	router.Get("/", HandleHealthz())
	router.Post("/signup", HandleSignup(logger, sessions, u))
	router.Post("/login", HandleLogin(logger, sessions, u))

	router.Group(func(r chi.Router) {
		r.Use(RequireAuthenticatedUser(logger, sessions, u))
		r.Post("/logout", HandleLogout(logger, sessions))

		r.Post("/create", HandleCreateTask(logger, t))
		r.Get("/all", HandleListTasks(logger, t))
		r.Get("/tasks/{task_id}", HandleGetTask(logger, t))
		r.Patch("/tasks/{task_id}/update", HandleUpdateTask(logger, t))
	})
	return router
}
