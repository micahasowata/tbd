package main

import (
	"net/http"

	"github.com/micahasowata/jason"
)

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}

	err := s.Read(w, r, &input)
	if err != nil {
		s.Write(w, http.StatusBadRequest, jason.Envelope{"error": err.Error()}, nil)
		return
	}

	err = s.validate.Struct(&input)
	if err != nil {
		s.Write(w, http.StatusUnprocessableEntity, jason.Envelope{"error": err.Error()}, nil)
		return
	}
	s.Write(w, http.StatusOK, jason.Envelope{"user": input}, nil)
}
