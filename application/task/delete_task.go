package task

import (
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type DeleteTask struct {
	repo domtask.Repository
}

func NewDeleteTask(repo domtask.Repository) *DeleteTask {
	return &DeleteTask{repo: repo}
}

func (uc *DeleteTask) Execute(id string) shared.Result[bool] {
	return uc.repo.Delete(id)
}
