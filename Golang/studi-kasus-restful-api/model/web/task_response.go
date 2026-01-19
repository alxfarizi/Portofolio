package web

import (
	"studi-kasus-restful-api/model/domain"
)

type TaskResponse struct {
	Id          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status"`
	Deadline    string            `json:"deadline,omitempty"`
}
