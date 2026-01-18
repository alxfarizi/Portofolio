package main

import (
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/repository"
	"studi-kasus-restful-api/service"
	"studi-kasus-restful-api/app"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db := app.NewDB()
	validate := validator.New()
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := httprouter.New()

	router.GET("/api/tasks", taskController.FindAll)
	router.GET("/api/tasks/:taskId", taskController.FindById)
	router.POST("/api/tasks", taskController.Create)
	router.PUT("/api/tasks/:taskId", taskController.Update)
	router.DELETE("/api/tasks/:taskId", taskController.Delete)
}
