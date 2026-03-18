package task

import (
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type UpdateTaskInput struct {
	ID       string
	Title    *string
	Status   *domtask.Status
	Priority *domtask.Priority
}

type UpdateTask struct {
	repo domtask.Repository
}

func NewUpdateTask(repo domtask.Repository) *UpdateTask {
	return &UpdateTask{repo: repo}
}

func (uc *UpdateTask) Execute(input UpdateTaskInput) shared.Result[*domtask.Task] {
	findResult := uc.repo.FindByID(input.ID)
	if findResult.IsErr() {
		return findResult
	}
	t := findResult.Value()

	if input.Title != nil {
		titleResult := t.UpdateTitle(*input.Title)
		if titleResult.IsErr() {
			return shared.Err[*domtask.Task](titleResult.Error())
		}
	}

	if input.Status != nil {
		transResult := t.TransitionTo(*input.Status)
		if transResult.IsErr() {
			return shared.Err[*domtask.Task](transResult.Error())
		}
	}

	if input.Priority != nil {
		prioResult := t.SetPriority(*input.Priority)
		if prioResult.IsErr() {
			return shared.Err[*domtask.Task](prioResult.Error())
		}
	}

	saveResult := uc.repo.Save(t)
	if saveResult.IsErr() {
		return shared.Err[*domtask.Task](saveResult.Error())
	}

	return shared.Ok(t)
}
