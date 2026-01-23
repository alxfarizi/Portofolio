package main

import (
	"fmt"
	"net/http"
	"studi-kasus-restful-api/middleware"

	"studi-kasus-restful-api/app"
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/repository"
	"studi-kasus-restful-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db := app.NewDB()
	validate := validator.New()
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := app.NewRouter(taskController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	fmt.Println("Server is running on http://localhost:3000")
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
