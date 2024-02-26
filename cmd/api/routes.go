package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/users/create", s.createUser)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", s.loginUser)

	return router
}
