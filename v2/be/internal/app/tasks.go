package app

import (
	"context"
	"errors"
	"net/http"

	"v2/be/internal/db"
	"v2/be/internal/models"
	"v2/be/internal/parser"
	"v2/be/internal/validator"

	"go.uber.org/zap"
)

type TaskCreater interface {
	Create(ctx context.Context, t *models.Task) error
}

func HandleCreateTask(logger *zap.Logger, tc TaskCreater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetUserID(r)

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			ReadError(w, logger, err)
			return
		}

		input.Title = parser.Sanitize(input.Title)
		input.Description = parser.Sanitize(input.Description)

		v := validator.New()
		v.RequiredString(input.Title, "title", validator.Required)
		v.RequiredString(input.Description, "description", validator.Required)
		if !v.Valid() {
			InvalidDataError(w, v.Errors())
			return
		}

		t := &models.Task{
			ID:          db.NewID(),
			UserID:      id,
			Title:       input.Title,
			Description: input.Description,
			Completed:   false,
		}

		err = tc.Create(r.Context(), t)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrDuplicateTask):
				DuplicateDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		err = parser.Write(w, http.StatusCreated, parser.Envelope{"payload": t.ID})
		if err != nil {
			writeError(w)
		}
	})
}

type TaskLister interface {
	All(ctx context.Context, userID string) ([]*models.Task, error)
}

func HandleListTasks(logger *zap.Logger, tl TaskLister) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetUserID(r)

		tasks, err := tl.All(r.Context(), id)
		if err != nil {
			ServerError(w, logger, err)
			return
		}

		if tasks == nil {
			tasks = []*models.Task{}
		}

		err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": tasks})
		if err != nil {
			writeError(w)
		}
	})
}

type TaskGetter interface {
	GetByID(ctx context.Context, id, userID string) (*models.Task, error)
}

func HandleGetTask(logger *zap.Logger, tg TaskGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetTaskID(r)
		userID := GetUserID(r)

		logger.Info(id)

		t, err := tg.GetByID(r.Context(), id, userID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				MissingDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		err = parser.Write(w, http.StatusFound, parser.Envelope{"payload": t})
		if err != nil {
			writeError(w)
		}
	})
}

type TaskUpdater interface {
	TaskGetter
	Update(ctx context.Context, t *models.Task) error
}

func HandleUpdateTask(logger *zap.Logger, tu TaskUpdater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetTaskID(r)
		userID := GetUserID(r)

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		err := parser.Read(w, r, &input)
		if err != nil {
			ReadError(w, logger, err)
			return
		}

		input.Title = parser.Sanitize(input.Title)
		input.Description = parser.Sanitize(input.Description)

		v := validator.New()
		v.RequiredString(input.Title, "title", validator.Required)
		v.RequiredString(input.Description, "description", validator.Required)
		if !v.Valid() {
			InvalidDataError(w, v.Errors())
			return
		}

		t, err := tu.GetByID(r.Context(), id, userID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				MissingDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		if t.Completed {
			err = parser.Write(w, http.StatusNotModified, parser.Envelope{"payload": "task completed"})
			if err != nil {
				writeError(w)
			}

			return
		}

		if t.Title != input.Title {
			t.Title = input.Title
		}

		if t.Description != input.Description {
			t.Description = input.Description
		}

		err = tu.Update(r.Context(), t)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				MissingDataError(w, logger, err)
			default:
				ServerError(w, logger, err)
			}
			return
		}

		err = parser.Write(w, http.StatusOK, parser.Envelope{"payload": t.ID})
		if err != nil {
			writeError(w)
		}
	})
}
