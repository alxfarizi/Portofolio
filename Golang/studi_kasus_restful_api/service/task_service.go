package service

import (
	"context"
	"studi-kasus-restful-api/model/web"
)

type TaskService interface {
	Create(ctx context.Context, request web.TaskCreateRequest) web.TaskResponse
	Update(ctx context.Context, request web.TaskUpdateRequest) web.TaskResponse
	Delete(ctx context.Context, taskId int)
	FindById(ctx context.Context, taskId int) web.TaskResponse
	FindAll(ctx context.Context) []web.TaskResponse
}
