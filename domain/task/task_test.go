package task_test

import (
	"testing"
	"time"

	"github.com/i-nishimura/goatodo/domain/task"
)

func TestNewTask(t *testing.T) {
	t.Run("valid title creates a task with todo status", func(t *testing.T) {
		result := task.NewTask("Buy milk")

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}

		tk := result.Value()
		if tk.Title() != "Buy milk" {
			t.Errorf("expected title 'Buy milk', got '%s'", tk.Title())
		}
		if tk.Status() != task.StatusTodo {
			t.Errorf("expected status todo, got %s", tk.Status())
		}
		if tk.Priority() != task.PriorityNone {
			t.Errorf("expected priority none, got %d", tk.Priority())
		}
		if tk.ID() == "" {
			t.Error("expected non-empty ID")
		}
	})

	t.Run("empty title returns error", func(t *testing.T) {
		result := task.NewTask("")

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
		if result.Error() != task.ErrEmptyTitle {
			t.Errorf("expected error '%s', got '%s'", task.ErrEmptyTitle, result.Error())
		}
	})
}

func TestTaskStatusTransition(t *testing.T) {
	t.Run("todo -> doing is valid", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.TransitionTo(task.StatusDoing)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Status() != task.StatusDoing {
			t.Errorf("expected status doing, got %s", tk.Status())
		}
	})

	t.Run("todo -> done is valid", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.TransitionTo(task.StatusDone)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Status() != task.StatusDone {
			t.Errorf("expected status done, got %s", tk.Status())
		}
	})

	t.Run("doing -> done is valid", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		tk.TransitionTo(task.StatusDoing)
		result := tk.TransitionTo(task.StatusDone)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Status() != task.StatusDone {
			t.Errorf("expected status done, got %s", tk.Status())
		}
	})

	t.Run("doing -> todo is valid (revert)", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		tk.TransitionTo(task.StatusDoing)
		result := tk.TransitionTo(task.StatusTodo)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Status() != task.StatusTodo {
			t.Errorf("expected status todo, got %s", tk.Status())
		}
	})

	t.Run("done -> todo is invalid", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		tk.TransitionTo(task.StatusDone)
		result := tk.TransitionTo(task.StatusTodo)

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
		if result.Error() != task.ErrInvalidTransition {
			t.Errorf("expected error '%s', got '%s'", task.ErrInvalidTransition, result.Error())
		}
	})

	t.Run("done -> doing is invalid", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		tk.TransitionTo(task.StatusDone)
		result := tk.TransitionTo(task.StatusDoing)

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})

	t.Run("same status transition is no-op (valid)", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.TransitionTo(task.StatusTodo)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
	})
}

func TestTaskSetPriority(t *testing.T) {
	t.Run("can set priority on a task", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.SetPriority(task.PriorityHigh)

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Priority() != task.PriorityHigh {
			t.Errorf("expected priority high, got %d", tk.Priority())
		}
	})

	t.Run("rejects invalid priority value", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.SetPriority(task.Priority(99))

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
		if result.Error() != task.ErrInvalidPriority {
			t.Errorf("expected error '%s', got '%s'", task.ErrInvalidPriority, result.Error())
		}
		if tk.Priority() != task.PriorityNone {
			t.Errorf("priority should remain unchanged, got %d", tk.Priority())
		}
	})

	t.Run("rejects negative priority value", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.SetPriority(task.Priority(-1))

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
	})
}

func TestTaskTransitionToInvalidStatus(t *testing.T) {
	t.Run("rejects unknown status string", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		result := tk.TransitionTo(task.Status("invalid"))

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
		if result.Error() != task.ErrInvalidStatus {
			t.Errorf("expected error '%s', got '%s'", task.ErrInvalidStatus, result.Error())
		}
	})
}

func TestTaskUpdateTitle(t *testing.T) {
	t.Run("can update title to valid value", func(t *testing.T) {
		tk := mustNewTask(t, "Old title")
		result := tk.UpdateTitle("New title")

		if result.IsErr() {
			t.Fatalf("expected Ok, got Err: %s", result.Error())
		}
		if tk.Title() != "New title" {
			t.Errorf("expected title 'New title', got '%s'", tk.Title())
		}
	})

	t.Run("cannot update title to empty", func(t *testing.T) {
		tk := mustNewTask(t, "Old title")
		result := tk.UpdateTitle("")

		if result.IsOk() {
			t.Fatal("expected Err, got Ok")
		}
		if tk.Title() != "Old title" {
			t.Errorf("title should remain unchanged, got '%s'", tk.Title())
		}
	})
}

func TestTaskSetDescription(t *testing.T) {
	t.Run("can set description", func(t *testing.T) {
		tk := mustNewTask(t, "Test task")
		tk.SetDescription("Some description")

		if tk.Description() != "Some description" {
			t.Errorf("expected description 'Some description', got '%s'", tk.Description())
		}
	})
}

func TestReconstruct(t *testing.T) {
	t.Run("reconstructs a task from persisted data", func(t *testing.T) {
		tk := task.Reconstruct("id-1", "Title", "Desc", task.StatusDoing, task.PriorityHigh, tk_time(), nil)

		if tk.ID() != "id-1" {
			t.Errorf("expected ID 'id-1', got '%s'", tk.ID())
		}
		if tk.Title() != "Title" {
			t.Errorf("expected title 'Title', got '%s'", tk.Title())
		}
		if tk.Description() != "Desc" {
			t.Errorf("expected description 'Desc', got '%s'", tk.Description())
		}
		if tk.Status() != task.StatusDoing {
			t.Errorf("expected status doing, got %s", tk.Status())
		}
		if tk.Priority() != task.PriorityHigh {
			t.Errorf("expected priority high, got %d", tk.Priority())
		}
		if tk.CompletedAt() != nil {
			t.Error("expected nil completedAt")
		}
	})
}

func tk_time() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func mustNewTask(t *testing.T, title string) *task.Task {
	t.Helper()
	result := task.NewTask(title)
	if result.IsErr() {
		t.Fatalf("failed to create task: %s", result.Error())
	}
	return result.Value()
}
