package task

import (
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type ListTasks struct {
	repo domtask.Repository
}

func NewListTasks(repo domtask.Repository) *ListTasks {
	return &ListTasks{repo: repo}
}

func (uc *ListTasks) Execute() shared.Result[[]*domtask.Task] {
	return uc.repo.FindAll()
}
