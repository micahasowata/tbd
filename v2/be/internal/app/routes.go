package app

import (
	"fmt"
	"net/http"
	"time"
	"v2/be/internal/models"
	"v2/be/internal/parser"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func Routes(
	sessions *scs.SessionManager,
	logger *zap.Logger,
	u *models.UsersModel,
	t *models.TasksModel,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: zap.NewStdLog(logger)}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(httprate.Limit(
		100,
		time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			err := parser.Write(w, http.StatusTooManyRequests, parser.Envelope{"error": fmt.Sprintf("rate limit exceed for %s. please retry later", r.URL.Path)})
			if err != nil {
				writeError(w)
			}
		}),
	))

	router.Use(sessions.LoadAndSave)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := parser.Write(w, http.StatusNotFound, parser.Envelope{"error": fmt.Sprintf("%s is not a valid endpoint", r.URL.Path)})
		if err != nil {
			writeError(w)
		}
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := parser.Write(w, http.StatusMethodNotAllowed, parser.Envelope{"error": fmt.Sprintf("%s is not supported for %s", r.Method, r.URL.Path)})
		if err != nil {
			writeError(w)
		}
	})

	router.Route("/v1", func(rr chi.Router) {
		rr.Group(func(ru chi.Router) {
			ru.Get("/", HandleHealthz())
			ru.Post("/signup", HandleSignup(logger, sessions, u))
			ru.Post("/login", HandleLogin(logger, sessions, u))
		})

		rr.Group(func(ra chi.Router) {
			ra.Use(RequireAuthenticatedUser(logger, sessions, u))
			ra.Post("/logout", HandleLogout(logger, sessions))

			ra.Post("/tasks/create", HandleCreateTask(logger, t))
			ra.Get("/tasks", HandleListTasks(logger, t))
			ra.Get("/tasks/{task_id}", HandleGetTask(logger, t))
			ra.Patch("/tasks/{task_id}/update", HandleUpdateTask(logger, t))
			ra.Patch("/tasks/{task_id}/complete", HandleCompleteTask(logger, t))
			ra.Delete("/tasks/{task_id}", HandleDeleteTask(logger, t))
		})

	})

	return cors.Default().Handler(router)
}
