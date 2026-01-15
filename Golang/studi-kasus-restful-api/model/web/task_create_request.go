package web

import (
	"studi-kasus-restful-api/model/domain"
	"time"
)

type TaskCreateRequest struct {
	Title       string            `validate:"required, min=1, max=255"`
	Description string            `validate:"required"`
	Status      domain.TaskStatus `default:"PENDING"`
	Deadline    time.Time         `validate:"required"`
}
