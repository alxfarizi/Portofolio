package web

import (
	"studi-kasus-restful-api/model/domain"
	"time"
)

type TaskCreateRequest struct {
	Title       string            `json:"title" validate:"required,min=1,max=255"`
	Description string            `json:"description" validate:"required"`
	Status      domain.TaskStatus `json:"status" validate:"required,oneof=pending in_progress completed"`
	Deadline    *time.Time        `json:"deadline" validate:"omitempty"`
}
