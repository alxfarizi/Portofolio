//go:build wireinject
// +build wireinject

package main

import (
	"net/http"
	"studi-kasus-restful-api/app"
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/middleware"
	"studi-kasus-restful-api/repository"
	"studi-kasus-restful-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

var taskSet = wire.NewSet(
	repository.NewTaskRepository,
	wire.Bind(new(repository.TasksRepository), new(*repository.TasksRepositoryImpl)),
	service.NewTaskService,
	wire.Bind(new(service.TaskService), new(*service.TaskServiceImpl)),
	controller.NewTaskController,
	wire.Bind(new(controller.TaskController), new(*controller.TaskControllerImpl)),
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewDB,
		NewValidator,
		taskSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)
	return nil
}
