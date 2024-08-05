package app

import (
	"net/http"
	"v2/be/internal/parser"

	"go.uber.org/zap"
)

func logError(logger *zap.Logger, err error) {
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
	}
}

func writeError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func ServerError(w http.ResponseWriter, logger *zap.Logger, err error) {
	logError(logger, err)

	perr := parser.Write(w, http.StatusInternalServerError, parser.Envelope{"error": "request could no longer be processed"})
	if perr != nil {
		writeError(w)
	}
}
