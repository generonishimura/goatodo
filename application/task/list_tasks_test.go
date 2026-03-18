package task_test

import (
	"testing"

	apptask "github.com/i-nishimura/goatodo/application/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestListTasks(t *testing.T) {
	t.Run("returns empty list when no tasks", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		uc := apptask.NewListTasks(repo)

		result := uc.Execute()

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if len(result.Value()) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(result.Value()))
		}
	})

	t.Run("returns all created tasks", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		createUC := apptask.NewCreateTask(repo)
		listUC := apptask.NewListTasks(repo)

		createUC.Execute("Task 1")
		createUC.Execute("Task 2")

		result := listUC.Execute()

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if len(result.Value()) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(result.Value()))
		}
	})
}
