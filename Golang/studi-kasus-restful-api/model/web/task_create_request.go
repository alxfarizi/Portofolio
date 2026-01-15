package web

import (
	"studi-kasus-restful-api/model/domain"
	"time"
)

type TaskCreateRequest struct {
	Title       string
	Description string
	Status      domain.TaskStatus
	Deadline    time.Time
}
