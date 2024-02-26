package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/micahasowata/jason"
)

func (s *server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		values := strings.Split(header, " ")
		if len(values) != 2 || values[0] != "Bearer" {
			s.Write(w, http.StatusUnauthorized, jason.Envelope{"error": "invalid authorization header"}, nil)
			return
		}

		token := []byte(values[1])

		claims, err := s.tokens.VerifyJWT(token)
		if err != nil {
			s.logger.Error(err.Error())
			s.Write(w, http.StatusForbidden, jason.Envelope{"error": "invalid token"}, nil)
			return
		}

		u, err := s.store.GetUserByEmail(claims.Email)
		if err != nil {
			s.logger.Error(err.Error())
			s.Write(w, http.StatusForbidden, jason.Envelope{"error": "token contains invalid data"}, nil)
			return
		}

		ctx := context.WithValue(r.Context(), userID, u.ID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
