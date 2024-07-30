package main

import (
	"net/http"
	"v2/be/internal/parser"

	"go.uber.org/zap"
)

func (app *application) logError(err error) {
	if err != nil {
		app.logger.Error(err.Error(), zap.Error(err))
	}
}

func (app *application) writeError(w http.ResponseWriter, err error) {
	app.logError(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (app *application) readError(w http.ResponseWriter, err error) {
	app.logError(err)

	err = parser.Write(w, http.StatusBadRequest, parser.Envelope{"error": err.Error()})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logError(err)

	err = parser.Write(w, http.StatusInternalServerError, parser.Envelope{"error": "request not processable"})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) dataConflictError(w http.ResponseWriter, err error) {
	app.logError(err)

	err = parser.Write(w, http.StatusConflict, parser.Envelope{"error": err.Error()})
	if err != nil {
		app.writeError(w, err)
	}
}
