package task_test

import (
	"testing"

	apptask "github.com/i-nishimura/goatodo/application/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestCreateTask(t *testing.T) {
	t.Run("creates a task with valid title", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		uc := apptask.NewCreateTask(repo)

		result := uc.Execute("Buy milk")

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		tk := result.Value()
		if tk.Title() != "Buy milk" {
			t.Errorf("expected title 'Buy milk', got '%s'", tk.Title())
		}

		// Verify persisted
		found := repo.FindByID(tk.ID())
		if found.IsErr() {
			t.Fatal("task was not persisted")
		}
	})

	t.Run("returns error for empty title", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		uc := apptask.NewCreateTask(repo)

		result := uc.Execute("")

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})
}
