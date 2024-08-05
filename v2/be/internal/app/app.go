package app

import (
	"net/http"
	"v2/be/internal/parser"

	"go.uber.org/zap"
)

func HandleHealthz(logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := parser.Write(w, http.StatusOK, parser.Envelope{"status": "OK"})
		if err != nil {
			writeError(w)
		}
	})
}
