package test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"studi-kasus-restful-api/app"
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/middleware"
	"studi-kasus-restful-api/repository"
	"studi-kasus-restful-api/service"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/studi_kasus_restful_api_test?parseTime=true")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := app.NewRouter(taskController)

	return middleware.NewAuthMiddleware(router)
}

func truncateTask(db *sql.DB) {
	db.Exec("TRUNCATE tasks")
}

func TestCreateTaskSuccess(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"Test Task","description":"Test Description","status":"done","deadline":"2026-01-20T00:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/tasks", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.NotNil(t, "Test Task", responseBody["data"].(map[string]interface{})["title"])
}

func TestCreateTaskFailed(t *testing.T) {

}

func TestUpdateTaskSuccess(t *testing.T) {

}

func TestUpdateTaskFailed(t *testing.T) {

}

func TestGetTaskSuccess(t *testing.T) {

}

func TestGetTaskFailed(t *testing.T) {

}

func TestDeleteTaskSuccess(t *testing.T) {

}

func TestDeleteTaskFailed(t *testing.T) {

}

func TestFindAllTaskSuccess(t *testing.T) {

}

func TestListTasksSuccess(t *testing.T) {

}

func TestListTasksFailed(t *testing.T) {

}
