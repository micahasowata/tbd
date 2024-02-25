package main

import "net/http"

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
