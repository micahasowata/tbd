package main

import (
	"net/http"
	"strconv"

	"github.com/micahasowata/jason"
)

func (s *server) createPost(w http.ResponseWriter, r *http.Request) {
	v, ok := r.Context().Value(userID).(int)
	if !ok {
		s.Write(w, http.StatusBadRequest, jason.Envelope{"error": "invalid request"}, nil)
		return
	}

	w.Write([]byte("id is" + strconv.Itoa(v)))
}
