package app

import (
	"context"
	"errors"
	"net/http"
	"v2/be/internal/db"
	"v2/be/internal/models"
	"v2/be/internal/parser"
	"v2/be/internal/validator"

	"github.com/alexedwards/argon2id"
	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

type UserCreater interface {
	Create(ctx context.Context, u *models.User) error
}

func HandleSignup(logger *zap.Logger, sessions *scs.SessionManager, uc UserCreater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			ReadError(w, logger, err)
			return
		}

		input.Username = parser.Sanitize(input.Username)
		input.Password = parser.Sanitize(input.Password)

		v := validator.New()
		v.RequiredString(input.Username, "username", validator.Required)
		v.RequiredString(input.Password, "password", validator.Required)
		v.MinString(input.Password, validator.MinPasswordLength, "password", "must be at least 8 characters")
		v.CheckPassword(input.Password, "password")

		if !v.Valid() {
			InvalidDataError(w, v.Errors())
			return
		}

		hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
		if err != nil {
			ServerError(w, logger, err)
			return
		}

		input.Password = hash

		u := &models.User{
			ID:       db.NewID(),
			Username: input.Username,
			Password: []byte(input.Password),
		}

		err = uc.Create(r.Context(), u)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrDuplicateUsername):
				DuplicateDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		sessions.Put(r.Context(), authenticatedUser, u.ID)

		err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": u.ID})
		if err != nil {
			writeError(w)
		}
	})
}

type UserGetter interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

func HandleLogin(logger *zap.Logger, sessions *scs.SessionManager, ug UserGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			ReadError(w, logger, err)
			return
		}

		input.Username = parser.Sanitize(input.Username)
		input.Password = parser.Sanitize(input.Password)

		v := validator.New()
		v.RequiredString(input.Username, "username", validator.Required)
		v.RequiredString(input.Password, "password", validator.Required)
		v.MinString(input.Password, validator.MinPasswordLength, "password", "must be at least 8 characters")
		v.CheckPassword(input.Password, "password")
		if !v.Valid() {
			InvalidDataError(w, v.Errors())
			return
		}

		u, err := ug.GetByUsername(r.Context(), input.Username)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				MissingDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		match, err := argon2id.ComparePasswordAndHash(input.Password, string(u.Password))
		if err != nil {
			ServerError(w, logger, err)
			return
		}

		if !match {
			MissingDataError(w, logger, models.ErrRecordNotFound)
			return
		}

		err = sessions.RenewToken(r.Context())
		if err != nil {
			ServerError(w, logger, err)
			return
		}

		sessions.Put(r.Context(), authenticatedUser, u.ID)

		err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": u.ID})
		if err != nil {
			writeError(w)
		}
	})
}
