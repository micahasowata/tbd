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

	err = parser.Write(w, http.StatusInternalServerError, parser.Envelope{"error": "request could no longer be processed"})
	if err != nil {
		writeError(w)
	}
}

func ReadError(w http.ResponseWriter, logger *zap.Logger, err error) {
	logError(logger, err)

	err = parser.Write(w, http.StatusBadRequest, parser.Envelope{"error": err.Error()})
	if err != nil {
		writeError(w)
	}
}

func InvalidDataError(w http.ResponseWriter, errs map[string]string) {
	err := parser.Write(w, http.StatusUnprocessableEntity, parser.Envelope{"error": errs})
	if err != nil {
		writeError(w)
	}
}

func DuplicateDataError(w http.ResponseWriter, logger *zap.Logger, err error) {
	logError(logger, err)

	err = parser.Write(w, http.StatusConflict, parser.Envelope{"error": err.Error()})
	if err != nil {
		writeError(w)
	}
}
