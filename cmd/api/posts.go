package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/micahasowata/jason"
	"github.com/micahasowata/tbd/pkg/domain"
	"github.com/micahasowata/tbd/pkg/store/sql/pg"
)

func (s *server) createPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userID).(int)
	if !ok {
		err := s.Write(w, http.StatusBadRequest, jason.Envelope{"error": "invalid request"}, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		return
	}

	var input struct {
		Title string `json:"title" validate:"required,min=1,max=100"`
		Body  string `json:"body" validate:"required"`
	}

	err := s.Read(w, r, &input)
	if err != nil {
		err := s.Write(w, http.StatusBadRequest, jason.Envelope{"error": err.Error()}, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		return
	}

	err = s.validate.Struct(&input)
	if err != nil {
		err := s.Write(w, http.StatusUnprocessableEntity, jason.Envelope{"error": err.Error()}, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
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
		err := s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "request could no longer be processed"}, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		return
	}

	err = s.Write(w, http.StatusOK, jason.Envelope{"post": post}, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
}

func (s *server) getPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userID).(int)
	if !ok {
		err := s.Write(w, http.StatusBadRequest, jason.Envelope{"error": "invalid request"}, nil)
		if err != nil {
			panic(err)
		}
		return
	}

	id := chi.URLParam(r, "id")

	postID, err := strconv.Atoi(id)
	if err != nil {
		err := s.Write(w, http.StatusUnprocessableEntity, jason.Envelope{"error": "invalid id"}, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		return
	}

	post := &domain.Post{
		UserID: userID,
		ID:     postID,
	}

	post, err = s.store.GetPost(post)
	if err != nil {
		switch {
		case errors.Is(err, pg.ErrPostNotFound):
			err := s.Write(w, http.StatusNotFound, jason.Envelope{"error": "post not found"}, nil)
			if err != nil {
				s.logger.Error(err.Error())
				return
			}
		default:
			err := s.Write(w, http.StatusInternalServerError, jason.Envelope{"error": "request could no longer be processed"}, nil)
			if err != nil {
				s.logger.Error(err.Error())
				return
			}
		}
		return
	}

	err = s.Write(w, http.StatusOK, jason.Envelope{"post": post}, nil)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
}
