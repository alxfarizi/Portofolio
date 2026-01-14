package repository

import (
	"context"
	"database/sql"
	"studi-kasus-restful-api/model/domain"
)

type TasksRepository interface {
	Save(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task
	Update(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task
	Delete(ctx context.Context, tx *sql.Tx, task domain.Task)
	FindById(ctx context.Context, tx *sql.Tx, taskId int) domain.Task
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Task
}
