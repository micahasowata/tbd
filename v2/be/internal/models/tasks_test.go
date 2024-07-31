package models

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestTasksModelCreate(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel2 := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Create", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(context.Background(), user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Test Task",
			Description: "This is a test task",
			Completed: true,
		}

		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)
	})

	t.Run("Duplicate Task", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(context.Background(), user)
		require.NoError(t, err)

		task1 := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Duplicate Task",
			Description: "This task will be duplicated",
		}

		err = tasksModel.Create(ctx, task1)
		require.NoError(t, err)

		task2 := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Duplicate Task",
			Description: "This is a duplicate task",
		}

		err = tasksModel.Create(ctx, task2)
		require.ErrorIs(t, err, ErrDuplicateTask)
	})

	t.Run("Empty Task", func(t *testing.T) {
		ctx := context.Background()
		err := tasksModel.Create(ctx, &Task{})
		require.Error(t, err)
	})

	t.Run("Empty Title", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(context.Background(), user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "",
			Description: "Task with empty title",
		}

		err = tasksModel.Create(ctx, task)
		require.Error(t, err)
	})

	t.Run("Empty Description", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(context.Background(), user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Task with empty description",
			Description: "",
		}

		err = tasksModel.Create(ctx, task)
		require.Error(t, err)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      uuid.New().String(),
			Title:       "Cancelled Context Task",
			Description: "This task should not be created",
		}

		err := tasksModel.Create(ctx, task)
		require.Error(t, err)
	})
}

func TestTasksModelAll(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel2 := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Retrieval", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			task := &Task{
				ID:          uuid.New().String(),
				UserID:      user.ID,
				Title:       fmt.Sprintf("Test Task %d", i+1),
				Description: fmt.Sprintf("This is test task %d", i+1),
			}
			err = tasksModel.Create(ctx, task)
			require.NoError(t, err)
		}

		tasks, err := tasksModel.All(ctx, user.ID)
		require.NoError(t, err)
		require.Len(t, tasks, 3)
		for _, task := range tasks {
			require.NotEmpty(t, task.ID)
			require.NotEmpty(t, task.Title)
			require.NotEmpty(t, task.Description)
			require.False(t, task.Completed)
		}
	})

	t.Run("No Tasks Found", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		tasks, err := tasksModel.All(ctx, user.ID)
		require.NoError(t, err)
		require.Empty(t, tasks)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		ctx := context.Background()

		tasks, err := tasksModel.All(ctx, "invalid_user_id")
		require.NoError(t, err)
		require.Empty(t, tasks)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tasks, err := tasksModel.All(ctx, uuid.New().String())
		require.Error(t, err)
		require.Empty(t, tasks)
	})
}

func TestTasksModelGetByID(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel2 := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Retrieval", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Test Task",
			Description: "This is a test task",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		retrievedTask, err := tasksModel.GetByID(ctx, task.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedTask)
		require.Equal(t, task.ID, retrievedTask.ID)
		require.Equal(t, task.Title, retrievedTask.Title)
		require.Equal(t, task.Description, retrievedTask.Description)
		require.Equal(t, task.Completed, retrievedTask.Completed)
	})

	t.Run("Task Not Found", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		nonExistentTaskID := uuid.New().String()
		retrievedTask, err := tasksModel.GetByID(ctx, nonExistentTaskID, user.ID)
		require.Error(t, err)
		require.Nil(t, retrievedTask)
		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Test Task",
			Description: "This is a test task",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		invalidUserID := uuid.New().String()
		retrievedTask, err := tasksModel.GetByID(ctx, task.ID, invalidUserID)
		require.Error(t, err)
		require.Nil(t, retrievedTask)
		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		retrievedTask, err := tasksModel.GetByID(ctx, uuid.New().String(), uuid.New().String())
		require.Error(t, err)
		require.Nil(t, retrievedTask)
	})
}

func TestTasksModelUpdate(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel2 := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Update", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Original Title",
			Description: "Original Description",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		updatedTask := &Task{
			ID:          task.ID,
			Title:       "Updated Title",
			Description: "Updated Description",
		}

		err = tasksModel.Update(ctx, updatedTask)
		require.NoError(t, err)

		retrievedTask, err := tasksModel.GetByID(ctx, task.ID, user.ID)
		require.NoError(t, err)
		require.Equal(t, updatedTask.Title, retrievedTask.Title)
		require.Equal(t, updatedTask.Description, retrievedTask.Description)
		require.False(t, retrievedTask.Completed)
	})

	t.Run("Update Completed Task", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Completed Task",
			Description: "This task is completed",
			Completed:   true,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		updatedTask := &Task{
			ID:          task.ID,
			Title:       "Updated Completed Task",
			Description: "This update should fail",
		}

		err = tasksModel.Update(ctx, updatedTask)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Update Non-existent Task", func(t *testing.T) {
		ctx := context.Background()

		nonExistentTask := &Task{
			ID:          uuid.New().String(),
			Title:       "Non-existent Task",
			Description: "This task doesn't exist",
		}

		err := tasksModel.Update(ctx, nonExistentTask)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Update to Duplicate Title", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel2.Create(ctx, user)
		require.NoError(t, err)

		task1 := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Existing Task",
			Description: "This is an existing task",
		}
		err = tasksModel.Create(ctx, task1)
		require.NoError(t, err)

		task2 := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Task to Update",
			Description: "This task will be updated",
		}
		err = tasksModel.Create(ctx, task2)
		require.NoError(t, err)

		updatedTask := &Task{
			ID:          task2.ID,
			Title:       "Existing Task",
			Description: "This update should fail due to duplicate title",
		}

		err = tasksModel.Update(ctx, updatedTask)
		require.ErrorIs(t, err, ErrDuplicateTask)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		task := &Task{
			ID:          uuid.New().String(),
			Title:       "Cancelled Context Task",
			Description: "This update should fail",
		}

		err := tasksModel.Update(ctx, task)
		require.Error(t, err)
	})
}

func TestTasksModelComplete(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Completion", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Test Task",
			Description: "This is a test task",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		err = tasksModel.Complete(ctx, task.ID, user.ID)
		require.NoError(t, err)

		completedTask, err := tasksModel.GetByID(ctx, task.ID, user.ID)
		require.NoError(t, err)
		require.True(t, completedTask.Completed)
	})

	t.Run("Complete Already Completed Task", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Already Completed Task",
			Description: "This task is already completed",
			Completed:   true,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		err = tasksModel.Complete(ctx, task.ID, user.ID)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Complete Non-existent Task", func(t *testing.T) {
		ctx := context.Background()

		nonExistentTaskID := uuid.New().String()
		userID := uuid.New().String()

		err := tasksModel.Complete(ctx, nonExistentTaskID, userID)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Complete Task with Wrong User ID", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user1 := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser1_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser2_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user2)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user1.ID,
			Title:       "User1's Task",
			Description: "This task belongs to user1",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		err = tasksModel.Complete(ctx, task.ID, user2.ID)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := tasksModel.Complete(ctx, uuid.New().String(), uuid.New().String())
		require.Error(t, err)
	})
}

func TestTasksModelDelete(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	usersModel := &UsersModel{pool: pool}
	tasksModel := &TasksModel{pool: pool}

	t.Run("Successful Delete", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user.ID,
			Title:       "Test Task",
			Description: "This is a test task",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		err = tasksModel.Delete(ctx, task.ID, user.ID)
		require.NoError(t, err)

		_, err = tasksModel.GetByID(ctx, task.ID, user.ID)
		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Delete Non-existent Task", func(t *testing.T) {
		ctx := context.Background()

		nonExistentTaskID := uuid.New().String()
		userID := uuid.New().String()

		err := tasksModel.Delete(ctx, nonExistentTaskID, userID)
		require.ErrorIs(t, err, ErrOpFailed)
	})

	t.Run("Delete Task with Wrong User ID", func(t *testing.T) {
		ctx := context.Background()

		hash, err := argon2id.CreateHash("Secret Password", argon2id.DefaultParams)
		require.NoError(t, err)

		user1 := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser1_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &User{
			ID:       uuid.New().String(),
			Username: fmt.Sprintf("testuser2_%s", time.Now().String()),
			Password: []byte(hash),
		}
		err = usersModel.Create(ctx, user2)
		require.NoError(t, err)

		task := &Task{
			ID:          uuid.New().String(),
			UserID:      user1.ID,
			Title:       "User1's Task",
			Description: "This task belongs to user1",
			Completed:   false,
		}
		err = tasksModel.Create(ctx, task)
		require.NoError(t, err)

		err = tasksModel.Delete(ctx, task.ID, user2.ID)
		require.ErrorIs(t, err, ErrOpFailed)

		// Verify the task still exists
		_, err = tasksModel.GetByID(ctx, task.ID, user1.ID)
		require.NoError(t, err)
	})

	t.Run("Cancelled Context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := tasksModel.Delete(ctx, uuid.New().String(), uuid.New().String())
		require.Error(t, err)
	})
}
