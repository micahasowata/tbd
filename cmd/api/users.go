package main

import (
	"net/http"
	"time"

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

func (s *server) loginUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
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

	user, err := s.store.GetUserByEmail(input.Email)
	if err != nil {
		s.Write(w, http.StatusForbidden, jason.Envelope{"error": "login failed; invalid email or password"}, nil)
		return
	}

	err = utils.CheckPassword(user.Password, input.Password)
	if err != nil {
		s.logger.Error(err.Error())
		s.Write(w, http.StatusForbidden, jason.Envelope{"error": "login failed; invalid email or password"}, nil)
		return
	}

	token, err := s.tokens.NewJWT(&domain.Claims{ID: user.ID, Email: user.Email}, 3*time.Hour)
	if err != nil {
		s.logger.Error(err.Error())
		s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "request could not be processed"}, nil)
		return
	}

	err = s.Write(w, http.StatusOK, jason.Envelope{"token": string(token)}, nil)
	if err != nil {
		return
	}
}
