package models_test

import (
	"context"
	"testing"
	"v2/be/internal/db"
	"v2/be/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestTasksCreate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.Phrase(),
			Description: gofakeit.AdjectiveQuantitative(),
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)
	})

	t.Run("duplicate", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		tOne := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       "run",
			Description: gofakeit.Phrase(),
		}

		err = tasks.Create(context.Background(), tOne)
		require.NoError(t, err)

		tTwo := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       "run",
			Description: gofakeit.Verb(),
		}

		err = tasks.Create(context.Background(), tTwo)
		require.NotNil(t, err)
		require.ErrorIs(t, err, models.ErrDuplicateTask)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      db.NewID(),
			Title:       gofakeit.HackerVerb(),
			Description: gofakeit.AppName(),
		}

		err := tasks.Create(ctx, task)
		require.Error(t, err)
	})
}

func TestTasksAll(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		for i := 0; i < 3; i++ {
			task := &models.Task{
				ID:          db.NewID(),
				UserID:      u.ID,
				Title:       gofakeit.VerbTransitive(),
				Description: gofakeit.Word(),
			}
			err = tasks.Create(context.Background(), task)
			require.NoError(t, err)
		}

		ts, err := tasks.All(context.Background(), u.ID)
		require.NoError(t, err)
		require.Len(t, ts, 3)

		for _, tt := range ts {
			require.NotEmpty(t, tt.ID)
			require.NotEmpty(t, tt.Title)
			require.NotEmpty(t, tt.Description)
			require.False(t, tt.Completed)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}

		tests := []struct {
			name string
			id   string
		}{
			{
				name: "valid user",
				id:   u.ID,
			},
			{
				name: "invalid user",
				id:   db.NewID(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ts, err := tasks.All(context.Background(), tt.id)
				require.NoError(t, err)
				require.Empty(t, ts)
			})
		}
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ts, err := tasks.All(ctx, db.NewID())
		require.Error(t, err)
		require.Empty(t, ts)
	})
}

func TestTasksGetByID(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.Verb(),
			Description: gofakeit.Blurb(),
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		rt, err := tasks.GetByID(context.Background(), task.ID, u.ID)
		require.NoError(t, err)
		require.NotNil(t, rt)
		require.Equal(t, task.ID, rt.ID)
		require.Equal(t, task.Title, rt.Title)
		require.Equal(t, task.Description, rt.Description)
		require.Equal(t, task.Completed, rt.Completed)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.Adverb(),
			Description: gofakeit.Phrase(),
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		tests := []struct {
			name   string
			id     string
			userID string
		}{
			{
				name:   "invalid task",
				id:     db.NewID(),
				userID: u.ID,
			},
			{
				name:   "invalid user",
				id:     task.ID,
				userID: db.NewID(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rt, err := tasks.GetByID(context.Background(), tt.id, tt.userID)
				require.Error(t, err)
				require.ErrorIs(t, err, models.ErrRecordNotFound)
				require.Nil(t, rt)
			})
		}
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		rt, err := tasks.GetByID(ctx, db.NewID(), db.NewID())
		require.Error(t, err)
		require.Nil(t, rt)
	})
}

func TestTasksUpdate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.AdjectiveDemonstrative(),
			Description: gofakeit.AdverbManner(),
			Completed:   false,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		ut := &models.Task{
			ID:          task.ID,
			Title:       "Updated Title",
			Description: "Updated Description",
		}

		err = tasks.Update(context.Background(), ut)
		require.NoError(t, err)

		rt, err := tasks.GetByID(context.Background(), task.ID, u.ID)
		require.NoError(t, err)
		require.Equal(t, ut.Title, rt.Title)
		require.Equal(t, ut.Description, rt.Description)
		require.False(t, rt.Completed)
	})

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.Animal(),
			Description: gofakeit.Question(),
			Completed:   true,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		ut := &models.Task{
			ID:          task.ID,
			Title:       gofakeit.AdjectiveIndefinite(),
			Description: gofakeit.Quote(),
		}

		err = tasks.Update(context.Background(), ut)
		require.ErrorIs(t, err, models.ErrOpFailed)
	})

	t.Run("invalid task", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			Title:       gofakeit.Adverb(),
			Description: gofakeit.Quote(),
		}

		err := tasks.Update(context.Background(), task)
		require.ErrorIs(t, err, models.ErrOpFailed)
	})

	t.Run("duplicate title", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		tOne := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.Adjective(),
			Description: gofakeit.Quote(),
		}

		err = tasks.Create(context.Background(), tOne)
		require.NoError(t, err)

		tTwo := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.VerbHelping(),
			Description: gofakeit.Animal(),
		}

		err = tasks.Create(context.Background(), tTwo)
		require.NoError(t, err)

		ut := &models.Task{
			ID:          tTwo.ID,
			Title:       tOne.Title,
			Description: gofakeit.AdjectiveProper(),
		}

		err = tasks.Update(context.Background(), ut)
		require.ErrorIs(t, err, models.ErrDuplicateTask)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		task := &models.Task{
			ID:          db.NewID(),
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Blurb(),
		}

		err := tasks.Update(ctx, task)
		require.Error(t, err)
	})
}

func TestTasksComplete(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Phrase(),
			Completed:   false,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		err = tasks.Complete(context.Background(), task.ID, u.ID)
		require.NoError(t, err)

		ct, err := tasks.GetByID(context.Background(), task.ID, u.ID)
		require.NoError(t, err)
		require.True(t, ct.Completed)
	})

	t.Run("completed task", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.AppName(),
			Description: gofakeit.Phrase(),
			Completed:   true,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		err = tasks.Complete(context.Background(), task.ID, u.ID)
		require.ErrorIs(t, err, models.ErrOpFailed)
	})

	t.Run("invalid task", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Phrase(),
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		tests := []struct {
			name   string
			id     string
			userID string
		}{
			{
				name:   "non-existent task",
				id:     db.NewID(),
				userID: u.ID,
			},
			{
				name:   "non-existent user",
				id:     task.ID,
				userID: db.NewID(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tasks.Complete(context.Background(), tt.id, tt.userID)
				require.ErrorIs(t, err, models.ErrOpFailed)

			})
		}
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := tasks.Complete(ctx, db.NewID(), db.NewID())
		require.Error(t, err)
	})
}

func TestTasksModelDelete(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Phrase(),
			Completed:   false,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		err = tasks.Delete(context.Background(), task.ID, u.ID)
		require.NoError(t, err)

		task, err = tasks.GetByID(context.Background(), task.ID, u.ID)
		require.ErrorIs(t, err, models.ErrRecordNotFound)
		require.Nil(t, task)
	})

	t.Run("non-existent task", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		users := &models.UsersModel{Pool: pool}
		u := &models.User{
			ID:       db.NewID(),
			Username: gofakeit.Username(),
			Password: []byte(testUserPassword(t)),
		}

		err := users.Create(context.Background(), u)
		require.NoError(t, err)

		tasks := &models.TasksModel{Pool: pool}
		task := &models.Task{
			ID:          db.NewID(),
			UserID:      u.ID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.Phrase(),
			Completed:   false,
		}

		err = tasks.Create(context.Background(), task)
		require.NoError(t, err)

		tests := []struct {
			name   string
			id     string
			userID string
		}{
			{
				name:   "invalid user",
				id:     task.ID,
				userID: db.NewID(),
			},
			{
				name:   "invalid task",
				id:     db.NewID(),
				userID: u.ID,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tasks.Delete(context.Background(), tt.id, tt.userID)
				require.ErrorIs(t, err, models.ErrOpFailed)
			})
		}
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		t.Parallel()

		pool := testPool(t)

		tasks := &models.TasksModel{Pool: pool}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := tasks.Delete(ctx, db.NewID(), db.NewID())
		require.Error(t, err)
	})
}
