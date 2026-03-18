package dto

import (
	"github.com/i-nishimura/goatodo/domain/task"
)

type TaskDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	CreatedAt   string `json:"createdAt"`
	CompletedAt string `json:"completedAt,omitempty"`
}

func FromTask(t *task.Task) TaskDTO {
	dto := TaskDTO{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Status:      string(t.Status()),
		Priority:    int(t.Priority()),
		CreatedAt:   t.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
	if t.CompletedAt() != nil {
		dto.CompletedAt = t.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
	}
	return dto
}

func FromTasks(tasks []*task.Task) []TaskDTO {
	dtos := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		dtos[i] = FromTask(t)
	}
	return dtos
}
