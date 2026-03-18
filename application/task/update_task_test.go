package task_test

import (
	"testing"

	apptask "github.com/i-nishimura/goatodo/application/task"
	domtask "github.com/i-nishimura/goatodo/domain/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestUpdateTask(t *testing.T) {
	t.Run("updates title and priority of existing task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		createUC := apptask.NewCreateTask(repo)
		updateUC := apptask.NewUpdateTask(repo)

		created := createUC.Execute("Old title")
		if created.IsErr() {
			t.Fatalf("setup failed: %s", created.Error())
		}

		newTitle := "New title"
		newPriority := domtask.PriorityHigh
		input := apptask.UpdateTaskInput{
			ID:       created.Value().ID(),
			Title:    &newTitle,
			Priority: &newPriority,
		}

		result := updateUC.Execute(input)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}

		found := repo.FindByID(created.Value().ID())
		if found.Value().Title() != "New title" {
			t.Errorf("expected title 'New title', got '%s'", found.Value().Title())
		}
		if found.Value().Priority() != domtask.PriorityHigh {
			t.Errorf("expected priority high, got %d", found.Value().Priority())
		}
	})

	t.Run("updates status transition", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		createUC := apptask.NewCreateTask(repo)
		updateUC := apptask.NewUpdateTask(repo)

		created := createUC.Execute("Test task")
		if created.IsErr() {
			t.Fatalf("setup failed: %s", created.Error())
		}

		newStatus := domtask.StatusDoing
		input := apptask.UpdateTaskInput{
			ID:     created.Value().ID(),
			Status: &newStatus,
		}

		result := updateUC.Execute(input)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}

		found := repo.FindByID(created.Value().ID())
		if found.Value().Status() != domtask.StatusDoing {
			t.Errorf("expected status doing, got %s", found.Value().Status())
		}
	})

	t.Run("returns error for non-existent task", func(t *testing.T) {
		repo := memory.NewTaskRepository()
		updateUC := apptask.NewUpdateTask(repo)

		input := apptask.UpdateTaskInput{ID: "non-existent"}
		result := updateUC.Execute(input)

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})
}
