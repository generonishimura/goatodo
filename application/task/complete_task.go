package task

import (
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type CompleteTask struct {
	repo domtask.Repository
}

func NewCompleteTask(repo domtask.Repository) *CompleteTask {
	return &CompleteTask{repo: repo}
}

func (uc *CompleteTask) Execute(id string) shared.Result[bool] {
	findResult := uc.repo.FindByID(id)
	if findResult.IsErr() {
		return shared.Err[bool](findResult.Error())
	}
	t := findResult.Value()
	transResult := t.TransitionTo(domtask.StatusDone)
	if transResult.IsErr() {
		return transResult
	}
	return uc.repo.Save(t)
}
