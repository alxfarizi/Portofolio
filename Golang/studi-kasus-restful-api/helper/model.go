package helper

import (
	"studi-kasus-restful-api/model/domain"
	"studi-kasus-restful-api/model/web"
)

func ToTaskResponse(task domain.Task) web.TaskResponse {
	var deadlineStr string
	if task.Deadline != nil {
		deadlineStr = task.Deadline.Format("2006-01-02")
	} else {
		deadlineStr = ""
	}

	return web.TaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Deadline:    deadlineStr,
	}
}

func ToTaskResponses(tasks []domain.Task) []web.TaskResponse {
	var taskResponses []web.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, ToTaskResponse(task))
	}

	return taskResponses
}
