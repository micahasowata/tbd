package main

import (
	"errors"
	"net/http"
	"v2/be/internal/db"
	"v2/be/internal/models"
	"v2/be/internal/parser"
	"v2/be/internal/validator"

	"github.com/alexedwards/argon2id"
	pv "github.com/wagslane/go-password-validator"
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

	v := validator.New()
	v.RequiredString(input.Username, "username", validator.Required)
	v.RequiredString(input.Password, "password", validator.Required)
	v.MinString(input.Password, validator.MinPasswordLength, "password", "must be at least 8 characters")
	err = pv.Validate(input.Password, validator.MinPasswordEntropy)
	if err != nil {
		v.AddError("password", err.Error())
	}

	if !v.Valid() {
		app.invalidDataError(w, v.Errors())
		return
	}

	hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		app.serverError(w, err)
		return
	}

	input.Password = hash

	u := &models.User{
		ID:       db.NewID(),
		Username: input.Username,
		Password: []byte(input.Password),
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

	app.sessions.Put(r.Context(), authenticatedUser, u.ID)

	err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": u.ID})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
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

	v := validator.New()
	v.RequiredString(input.Username, "username", validator.Required)
	v.RequiredString(input.Password, "password", validator.Required)
	v.MinString(input.Password, validator.MinPasswordLength, "password", "must be at least 8 characters")
	err = pv.Validate(input.Password, validator.MinPasswordEntropy)
	if err != nil {
		v.AddError("password", err.Error())
	}

	if !v.Valid() {
		app.invalidDataError(w, v.Errors())
		return
	}

	u, err := app.models.Users.GetByUsername(r.Context(), input.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, string(u.Password))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !match {
		app.recordNotFoundError(w, models.ErrRecordNotFound)
		return
	}

	err = app.sessions.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessions.Put(r.Context(), authenticatedUser, u.ID)

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": u.ID})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	id := getIDFromCtx(r)

	err := app.sessions.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessions.Remove(r.Context(), authenticatedUser)

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": id})
	if err != nil {
		app.writeError(w, err)
	}
}
