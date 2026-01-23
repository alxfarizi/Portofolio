package app

import (
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(taskController controller.TaskController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/tasks", taskController.FindAll)
	router.GET("/api/tasks/:taskId", taskController.FindById)
	router.POST("/api/tasks", taskController.Create)
	router.PUT("/api/tasks/:taskId", taskController.Update)
	router.DELETE("/api/tasks/:taskId", taskController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}
