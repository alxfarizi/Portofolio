package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"studi-kasus-restful-api/app"
	"studi-kasus-restful-api/controller"
	"studi-kasus-restful-api/helper"
	"studi-kasus-restful-api/middleware"
	"studi-kasus-restful-api/model/domain"
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

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, "Test Task", data["title"])
	assert.Equal(t, "Test Description", data["description"])
	assert.Equal(t, "done", data["status"])
	assert.Equal(t, "2026-01-20", data["deadline"])
}

func TestCreateTaskFailed(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"","description":"Test Description","status":"done","deadline":"2026-01-20T00:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/tasks", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateTaskSuccess(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusDone,
		Deadline:    nil,
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"Test Task Updated","description":"Test Description Updated","status":"in_progress","deadline":"2026-02-20T00:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/tasks/"+strconv.Itoa(task.Id), requestBody)
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

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, task.Id, int(data["id"].(float64)))
	assert.Equal(t, "Test Task Updated", data["title"])
	assert.Equal(t, "Test Description Updated", data["description"])
	assert.Equal(t, "in_progress", data["status"])
	assert.Equal(t, "2026-02-20", data["deadline"])
}

func TestUpdateTaskFailed(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusDone,
		Deadline:    nil,
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"","description":"Test Description","status":"in_progress","deadline":"2026-02-20T00:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/tasks/"+strconv.Itoa(task.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateTaskNotFound(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title":"Test Task Updated","description":"Test Description Updated","status":"in_progress","deadline":"2026-02-20T00:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/tasks/9999", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestGetTaskSuccess(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusDone,
		Deadline:    nil,
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/tasks/"+strconv.Itoa(task.Id), nil)
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

	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, task.Id, int(data["id"].(float64)))

	assert.Equal(t, task.Title, data["title"])
	assert.Equal(t, task.Description, data["description"])
	assert.Equal(t, string(task.Status), data["status"])

	assert.Nil(t, data["deadline"])
}

func TestGetTaskFailed(t *testing.T) {

	db := setupTestDB()
	truncateTask(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/tasks/404", nil)
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteTaskSuccess(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusDone,
		Deadline:    nil,
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/tasks/"+strconv.Itoa(task.Id), nil)
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
}

func TestDeleteTaskFailed(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/tasks/404", nil)
	request.Header.Add("TASK-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestFindAllTaskSuccess(t *testing.T) {
	db := setupTestDB()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	// Insert task 1
	task1 := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task 1",
		Description: "Test Description 1",
		Status:      domain.StatusPending,
	})
	// Insert task 2
	task2 := taskRepository.Save(context.Background(), tx, domain.Task{
		Title:       "Test Task 2",
		Description: "Test Description 2",
		Status:      domain.StatusDone,
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/tasks", nil)
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

	tasks := responseBody["data"].([]interface{})
	assert.Equal(t, 2, len(tasks))

	task1Data := tasks[0].(map[string]interface{})
	task2Data := tasks[1].(map[string]interface{})

	assert.Equal(t, task1.Id, int(task1Data["id"].(float64)))
	assert.Equal(t, task1.Title, task1Data["title"])
	assert.Equal(t, task1.Description, task1Data["description"])
	assert.Equal(t, string(task1.Status), task1Data["status"])

	assert.Equal(t, task2.Id, int(task2Data["id"].(float64)))
	assert.Equal(t, task2.Title, task2Data["title"])
	assert.Equal(t, task2.Description, task2Data["description"])
	assert.Equal(t, string(task2.Status), task2Data["status"])
}
