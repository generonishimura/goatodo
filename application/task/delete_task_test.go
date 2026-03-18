package task_test

import (
	"testing"

	apptask "github.com/i-nishimura/goatodo/application/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestDeleteTask(t *testing.T) {
	t.Run("deletes an existing task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		createUC := apptask.NewCreateTask(repo)
		deleteUC := apptask.NewDeleteTask(repo)

		created := createUC.Execute("Test task")
		if created.IsErr() {
			t.Fatalf("setup failed: %s", created.Error())
		}

		result := deleteUC.Execute(created.Value().ID())

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}

		// Verify deleted
		found := repo.FindByID(created.Value().ID())
		if found.IsOk() {
			t.Error("task should have been deleted")
		}
	})

	t.Run("returns error for non-existent task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		uc := apptask.NewDeleteTask(repo)

		result := uc.Execute("non-existent-id")

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})
}
