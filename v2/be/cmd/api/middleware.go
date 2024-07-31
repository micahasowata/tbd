package main

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrUnauthenticatedUser = errors.New("requires authenticated user")
)

type ctx string

var userID = ctx("userID")

func getIDFromCtx(r *http.Request) string {
	return r.Context().Value(userID).(string)
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := app.sessions.Exists(r.Context(), authenticatedUser)
		if !ok {
			app.unauthorizedAccessError(w, ErrUnauthenticatedUser)
			return
		}

		id := app.sessions.GetString(r.Context(), authenticatedUser)
		if len(id) == 0 {
			app.unauthorizedAccessError(w, ErrUnauthenticatedUser)
			return
		}

		exists, err := app.models.Users.Exists(r.Context(), id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !exists {
			app.unauthorizedAccessError(w, ErrUnauthenticatedUser)
			return
		}

		ctx := context.WithValue(r.Context(), userID, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
