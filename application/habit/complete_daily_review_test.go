package habit

import (
	"testing"
	"time"

	domhabit "github.com/i-nishimura/goatodo/domain/habit"
	domtask "github.com/i-nishimura/goatodo/domain/task"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestCompleteDailyReview_Execute(t *testing.T) {
	t.Run("creates and completes a daily review for today", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		taskRepo := memory.NewTaskRepository()
		uc := NewCompleteDailyReview(repo, taskRepo)

		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		result := uc.Execute(today)

		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		review := result.Value()
		if review.Status() != domhabit.ReviewCompleted {
			t.Errorf("expected status %s, got %s", domhabit.ReviewCompleted, review.Status())
		}
		if review.CompletedAt() == nil {
			t.Error("expected completedAt to be set")
		}
	})

	t.Run("records task counts from task repository", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		taskRepo := memory.NewTaskRepository()

		// Create some tasks: 2 done, 1 todo
		createAndSaveTask(t, taskRepo, "task1", "done")
		createAndSaveTask(t, taskRepo, "task2", "done")
		createAndSaveTask(t, taskRepo, "task3", "todo")

		uc := NewCompleteDailyReview(repo, taskRepo)
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		result := uc.Execute(today)

		review := result.Value()
		if review.CompletedTaskCount() != 2 {
			t.Errorf("expected 2 completed tasks, got %d", review.CompletedTaskCount())
		}
		if review.TotalTaskCount() != 3 {
			t.Errorf("expected 3 total tasks, got %d", review.TotalTaskCount())
		}
	})

	t.Run("updates existing pending review instead of creating duplicate", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		taskRepo := memory.NewTaskRepository()

		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		existing := domhabit.NewDailyReview(today).Value()
		repo.Save(existing)

		uc := NewCompleteDailyReview(repo, taskRepo)
		result := uc.Execute(today)

		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		review := result.Value()
		if review.ID() != existing.ID() {
			t.Error("expected to reuse existing review, not create new one")
		}
	})

	t.Run("returns error if review already completed", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		taskRepo := memory.NewTaskRepository()

		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		existing := domhabit.NewDailyReview(today).Value()
		existing.RecordTaskCounts(1, 1)
		existing.Complete()
		repo.Save(existing)

		uc := NewCompleteDailyReview(repo, taskRepo)
		result := uc.Execute(today)

		if result.IsOk() {
			t.Error("expected error for already completed review")
		}
	})
}

func createAndSaveTask(t *testing.T, repo *memory.TaskRepository, title, status string) {
	t.Helper()
	taskResult := domtask.NewTask(title)
	if taskResult.IsErr() {
		t.Fatal(taskResult.Error())
	}
	task := taskResult.Value()
	if status == "done" {
		task.TransitionTo(domtask.StatusDone)
	}
	repo.Save(task)
}
