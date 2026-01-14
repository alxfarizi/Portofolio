package domain

import "time"

type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
)

type Task struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Deadline    time.Time  `json:"deadline"`
}
