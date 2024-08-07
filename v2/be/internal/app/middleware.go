package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

type CtxKey string

type UserExister interface {
	Exists(ctx context.Context, id string) (bool, error)
}

const authenticatedUser = "authenticatedUser"

var (
	ErrUnauthorized = errors.New("must be an authenticated user")
	userID          = CtxKey("userID")
)

// RequireAuthenticatedUser returns a function that satisfies the chi middleware pattern
func RequireAuthenticatedUser(logger *zap.Logger, sessions *scs.SessionManager, ue UserExister) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := sessions.Exists(r.Context(), authenticatedUser)
			if !ok {
				UnauthorizedAccessError(w, logger, ErrUnauthorized)
				return
			}

			id := strings.TrimSpace(sessions.GetString(r.Context(), authenticatedUser))
			if len(id) == 0 {
				UnauthorizedAccessError(w, logger, ErrUnauthorized)
				return
			}

			exists, err := ue.Exists(r.Context(), id)
			if err != nil {
				ServerError(w, logger, err)
				return
			}

			if !exists {
				UnauthorizedAccessError(w, logger, ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userID, id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) string {
	return r.Context().Value(userID).(string)
}
