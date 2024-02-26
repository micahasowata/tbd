package main

import (
	"net/http"

	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
)

func (s *server) createPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userID).(int)
	if !ok {
		s.Write(w, http.StatusBadRequest, jason.Envelope{"error": "invalid request"}, nil)
		return
	}

	var input struct {
		Title string `json:"title" validate:"required,min=1,max=100"`
		Body  string `json:"body" validate:"required"`
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

	post := &domain.Post{
		UserID: userID,
		Title:  input.Title,
		Body:   input.Body,
	}

	post, err = s.store.CreatePost(post)
	if err != nil {
		s.logger.Error(err.Error())
		s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "request could no longer be processed"}, nil)
		return
	}

	s.Write(w, http.StatusOK, jason.Envelope{"post": post}, nil)
}
