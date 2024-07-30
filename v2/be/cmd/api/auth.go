package main

import (
	"errors"
	"net/http"
	"v2/be/internal/db"
	"v2/be/internal/models"
	"v2/be/internal/parser"

	"github.com/alexedwards/argon2id"
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

	hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		app.serverError(w, err)
		return
	}

	input.Password = hash

	u := models.User{
		ID:       db.NewID(),
		Username: input.Username,
		Password: input.Password,
	}

	err = app.models.Users.Create(r.Context(), u)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateUsername):
			app.dataConflictError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": u.ID})
	if err != nil {
		app.writeError(w, err)
	}
}
