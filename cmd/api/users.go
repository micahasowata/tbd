package main

import (
	"net/http"

	"github.com/micahasowata/jason"
)

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := s.Read(w, r, &input)
	if err != nil {
		s.Write(w, http.StatusBadRequest, jason.Envelope{"error": err.Error()}, nil)
		return
	}

	s.Write(w, http.StatusOK, jason.Envelope{"user": input}, nil)
}
