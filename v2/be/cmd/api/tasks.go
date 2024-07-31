package main

import (
	"errors"
	"net/http"
	"v2/be/internal/db"
	"v2/be/internal/models"
	"v2/be/internal/parser"
	"v2/be/internal/validator"
)

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromCtx(r)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	err := parser.Read(w, r, &input)
	if err != nil {
		app.readError(w, err)
		return
	}

	input.Title = parser.Sanitize(input.Title)
	input.Description = parser.Sanitize(input.Description)

	v := validator.New()
	v.RequiredString(input.Title, "title", validator.Required)
	v.RequiredString(input.Description, "description", validator.Required)
	if !v.Valid() {
		app.invalidDataError(w, v.Errors())
		return
	}

	t := &models.Task{
		ID:          db.NewID(),
		UserID:      id,
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	err = app.models.Tasks.Create(r.Context(), t)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateTask):
			app.dataConflictError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": t.ID})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) allTasks(w http.ResponseWriter, r *http.Request) {
	id := getIDFromCtx(r)
	tasks, err := app.models.Tasks.All(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if tasks == nil {
		tasks = []*models.Task{}
	}

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": tasks})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) getTask(w http.ResponseWriter, r *http.Request) {
	id := getTaskID(r)
	userID := getIDFromCtx(r)

	t, err := app.models.Tasks.GetByID(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	err = parser.Write(w, http.StatusFound, parser.Envelope{"payload": t})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	id := getTaskID(r)
	userID := getIDFromCtx(r)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	err := parser.Read(w, r, &input)
	if err != nil {
		app.readError(w, err)
		return
	}

	input.Title = parser.Sanitize(input.Title)
	input.Description = parser.Sanitize(input.Description)

	v := validator.New()
	v.RequiredString(input.Title, "title", validator.Required)
	v.RequiredString(input.Description, "description", validator.Required)
	if !v.Valid() {
		app.invalidDataError(w, v.Errors())
		return
	}

	t, err := app.models.Tasks.GetByID(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	if t.Completed {
		err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": "task completed"})
		if err != nil {
			app.writeError(w, err)
		}

		return
	}

	if t.Title != input.Title {
		t.Title = input.Title
	}

	if t.Description != input.Description {
		t.Description = input.Description
	}

	err = app.models.Tasks.Update(r.Context(), t)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": t.ID})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) completeTask(w http.ResponseWriter, r *http.Request) {
	id := getTaskID(r)
	userID := getIDFromCtx(r)

	t, err := app.models.Tasks.GetByID(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	if t.Completed {
		err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": "task completed"})
		if err != nil {
			app.writeError(w, err)
		}
		return
	}

	err = app.models.Tasks.Complete(r.Context(), id, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": t.ID})
	if err != nil {
		app.writeError(w, err)
	}
}

func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := getTaskID(r)
	userID := getIDFromCtx(r)

	t, err := app.models.Tasks.GetByID(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.recordNotFoundError(w, err)
		default:
			app.serverError(w, err)
		}
		return
	}

	err = app.models.Tasks.Delete(r.Context(), t.ID, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": "deleted"})
	if err != nil {
		app.writeError(w, err)
	}
}
