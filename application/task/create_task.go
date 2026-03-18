package task

import (
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type CreateTask struct {
	repo domtask.Repository
}

func NewCreateTask(repo domtask.Repository) *CreateTask {
	return &CreateTask{repo: repo}
}

func (uc *CreateTask) Execute(title string) shared.Result[*domtask.Task] {
	result := domtask.NewTask(title)
	if result.IsErr() {
		return result
	}
	saveResult := uc.repo.Save(result.Value())
	if saveResult.IsErr() {
		return shared.Err[*domtask.Task](saveResult.Error())
	}
	return result
}
