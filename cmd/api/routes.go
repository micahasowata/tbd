package main

import (
	"net/http"
)

func (s *server) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /v1/users/create", s.createUser)
	return router
}
