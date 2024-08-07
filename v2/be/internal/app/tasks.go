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
