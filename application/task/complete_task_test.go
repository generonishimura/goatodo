package task_test

import (
	"testing"

	apptask "github.com/i-nishimura/goatodo/application/task"
	domtask "github.com/i-nishimura/goatodo/domain/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestCompleteTask(t *testing.T) {
	t.Run("completes an existing task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		createUC := apptask.NewCreateTask(repo)
		completeUC := apptask.NewCompleteTask(repo)

		created := createUC.Execute("Test task")
		if created.IsErr() {
			t.Fatalf("setup failed: %s", created.Error())
		}

		result := completeUC.Execute(created.Value().ID())

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}

		found := repo.FindByID(created.Value().ID())
		if found.Value().Status() != domtask.StatusDone {
			t.Errorf("expected status done, got %s", found.Value().Status())
		}
	})

	t.Run("returns error for non-existent task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		uc := apptask.NewCompleteTask(repo)

		result := uc.Execute("non-existent-id")

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})
}
