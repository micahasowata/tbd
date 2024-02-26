package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *server) routes() http.Handler {
	router := chi.NewRouter()

	router.Post("/v1/users/create", s.createUser)
	router.Post("/v1/users/login", s.loginUser)

	router.With(s.authUser).Post("/v1/posts/create", s.createPost)

	return router
}
