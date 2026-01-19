package main

import (
	"fmt"
	"net/http"
	"studi-kasus-restful-api/exception"
	"studi-kasus-restful-api/middleware"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"studi-kasus-restful-api/app"
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/repository"
	"studi-kasus-restful-api/service"
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

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	fmt.Println("Server is running on http://localhost:3000")
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
