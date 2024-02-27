package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"
)

func (s *server) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  zap.NewStdLog(s.logger),
		NoColor: true,
	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.CleanPath)
	router.Use(middleware.RealIP)
	router.Use(httprate.Limit(4, 1*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint)))

	router.Post("/v1/users/create", s.createUser)
	router.Post("/v1/users/login", s.loginUser)

	router.With(s.authUser).Post("/v1/posts/create", s.createPost)
	router.With(s.authUser).Get("/v1/posts/{id}", s.getPost)

	return router
}
