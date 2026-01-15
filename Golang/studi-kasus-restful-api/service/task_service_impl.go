package service

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/model/domain"
	"studi-kasus-restful-api/model/web"
	"studi-kasus-restful-api/repository"
)

type TaskServiceImpl struct {
	TaskRepository repository.TasksRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func (service *TaskServiceImpl) Create(ctx context.Context, request web.TaskCreateRequest) web.TaskResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOdRollback(tx)

	task := domain.Task{
		Title:       request.Title,
		Description: request.Description,
		Deadline:    request.Deadline,
		Status:      request.Status,
	}

	task = service.TaskRepository.Save(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

func (service *TaskServiceImpl) Update(ctx context.Context, request web.TaskUpdateRequest) web.TaskResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOdRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, request.Id)
	helper.PanicIfError(err)

	task.Title = request.Title
	task.Description = request.Description
	task.Status = request.Status
	task.Deadline = request.Deadline

	task = service.TaskRepository.Update(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

func (service *TaskServiceImpl) Delete(ctx context.Context, taskId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOdRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, taskId)
	helper.PanicIfError(err)

	service.TaskRepository.Delete(ctx, tx, task)
}

func (service *TaskServiceImpl) FindById(ctx context.Context, taskId int) web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOdRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, taskId)
	helper.PanicIfError(err)

	return helper.ToTaskResponse(task)
}

func (service *TaskServiceImpl) FindAll(ctx context.Context) []web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOdRollback(tx)

	tasks := service.TaskRepository.FindAll(ctx, tx)

	return helper.ToTaskResponses(tasks)
}
