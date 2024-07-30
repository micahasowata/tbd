package main

import (
	"net/http"
	"v2/be/internal/parser"
)

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := parser.Read(w, r, &input)
	if err != nil {
		app.readError(w, err)
		return
	}

	input.Username = parser.Sanitize(input.Username)
	input.Password = parser.Sanitize(input.Password)

	err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": input})
	if err != nil {
		app.writeError(w, err)
	}
}
