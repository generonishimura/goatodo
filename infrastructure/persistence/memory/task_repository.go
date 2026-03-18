package memory

import (
	"sort"
	"sync"

	"github.com/i-nishimura/goatodo/domain/shared"
	"github.com/i-nishimura/goatodo/domain/task"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*task.Task),
	}
}

func (r *TaskRepository) Save(t *task.Task) shared.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID()] = t
	return shared.Ok(true)
}

func (r *TaskRepository) FindByID(id string) shared.Result[*task.Task] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return shared.Err[*task.Task]("task not found")
	}
	return shared.Ok(t)
}

func (r *TaskRepository) FindAll() shared.Result[[]*task.Task] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*task.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt().After(result[j].CreatedAt())
	})
	return shared.Ok(result)
}

func (r *TaskRepository) Delete(id string) shared.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return shared.Err[bool]("task not found")
	}
	delete(r.tasks, id)
	return shared.Ok(true)
}
