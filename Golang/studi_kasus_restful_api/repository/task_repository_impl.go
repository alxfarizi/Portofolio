package repository

import (
	"context"
	"database/sql"
	"errors"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/model/domain"
)

type TasksRepositoryImpl struct {
}

func NewTaskRepository() TasksRepository {
	return &TasksRepositoryImpl{}
}

func (repository *TasksRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	SQL := "INSERT INTO tasks (title, description, deadline, status) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, task.Title, task.Description, task.Deadline, task.Status)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	task.Id = int(id)
	return task
}

func (repository *TasksRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	SQL := "UPDATE tasks SET title = ?, description = ?, deadline = ?, status = ? WHERE id = ?"

	var deadline interface{}
	if task.Deadline != nil && !task.Deadline.IsZero() {
		deadline = task.Deadline
	} else {
		deadline = nil
	}

	_, err := tx.ExecContext(ctx, SQL, task.Title, task.Description, deadline, task.Status, task.Id)
	helper.PanicIfError(err)

	return task
}

func (repository *TasksRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, task domain.Task) {
	SQL := "DELETE FROM tasks WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, task.Id)
	helper.PanicIfError(err)
}

func (repository *TasksRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, taskId int) (domain.Task, error) {
	SQL := "SELECT id, title, description, deadline, status FROM tasks WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, taskId)
	helper.PanicIfError(err)
	defer rows.Close()

	task := domain.Task{}

	if rows.Next() {
		var deadline sql.NullTime
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &deadline, &task.Status)
		helper.PanicIfError(err)

		if deadline.Valid {
			task.Deadline = &deadline.Time
		} else {
			task.Deadline = nil
		}

		return task, nil
	} else {
		return task, errors.New("task is not found")
	}
}

func (repository *TasksRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Task {
	SQL := "SELECT id, title, description, deadline, status FROM tasks"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		task := domain.Task{}
		var deadline sql.NullTime
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &deadline, &task.Status)
		helper.PanicIfError(err)

		if deadline.Valid {
			task.Deadline = &deadline.Time
		} else {
			task.Deadline = nil
		}

		tasks = append(tasks, task)
	}

	return tasks
}
