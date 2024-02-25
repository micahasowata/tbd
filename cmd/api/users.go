package main

import (
	"net/http"

	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/utils"
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

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		s.logger.Error(err.Error())
		s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "could not process request"}, nil)
		return
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hash,
	}

	user, err = s.store.CreateUser(user)
	if err != nil {
		s.logger.Error(err.Error())
		s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "could not process request"}, nil)
		return
	}

	s.Write(w, http.StatusOK, jason.Envelope{"user": user}, nil)
}
