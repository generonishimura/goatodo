package task

import "github.com/i-nishimura/goatodo/domain/shared"

type Repository interface {
	Save(task *Task) shared.Result[bool]
	FindByID(id string) shared.Result[*Task]
	FindAll() shared.Result[[]*Task]
	Delete(id string) shared.Result[bool]
}
