package web

import (
	"studi-kasus-restful-api/model/domain"
	"time"
)

type TaskResponse struct {
	Id          int
	Title       string
	Description string
	Status      domain.TaskStatus
	Deadline    time.Time
}
