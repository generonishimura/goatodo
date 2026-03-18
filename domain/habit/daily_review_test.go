package habit

import (
	"testing"
	"time"
)

func TestNewDailyReview(t *testing.T) {
	t.Run("creates a daily review for a given date", func(t *testing.T) {
		date := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		result := NewDailyReview(date)

		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		review := result.Value()
		if review.ID() == "" {
			t.Error("expected non-empty ID")
		}
		if !review.Date().Equal(date) {
			t.Errorf("expected date %v, got %v", date, review.Date())
		}
		if review.Status() != ReviewPending {
			t.Errorf("expected status %s, got %s", ReviewPending, review.Status())
		}
		if review.CompletedTaskCount() != 0 {
			t.Errorf("expected 0 completed tasks, got %d", review.CompletedTaskCount())
		}
		if review.TotalTaskCount() != 0 {
			t.Errorf("expected 0 total tasks, got %d", review.TotalTaskCount())
		}
	})

	t.Run("normalizes date to midnight UTC", func(t *testing.T) {
		date := time.Date(2026, 3, 18, 15, 30, 0, 0, time.UTC)
		result := NewDailyReview(date)

		review := result.Value()
		expected := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		if !review.Date().Equal(expected) {
			t.Errorf("expected date %v, got %v", expected, review.Date())
		}
	})
}

func TestDailyReview_RecordTaskCounts(t *testing.T) {
	t.Run("records completed and total task counts", func(t *testing.T) {
		date := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		review := NewDailyReview(date).Value()

		result := review.RecordTaskCounts(3, 5)
		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		if review.CompletedTaskCount() != 3 {
			t.Errorf("expected 3 completed tasks, got %d", review.CompletedTaskCount())
		}
		if review.TotalTaskCount() != 5 {
			t.Errorf("expected 5 total tasks, got %d", review.TotalTaskCount())
		}
	})

	t.Run("rejects negative completed count", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()
		result := review.RecordTaskCounts(-1, 5)
		if result.IsOk() {
			t.Error("expected error for negative completed count")
		}
		if result.Error() != ErrInvalidTaskCount {
			t.Errorf("expected %q, got %q", ErrInvalidTaskCount, result.Error())
		}
	})

	t.Run("rejects completed count exceeding total", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()
		result := review.RecordTaskCounts(6, 5)
		if result.IsOk() {
			t.Error("expected error for completed > total")
		}
	})
}

func TestDailyReview_Complete(t *testing.T) {
	t.Run("completes a pending review", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()
		review.RecordTaskCounts(3, 5)

		result := review.Complete()
		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		if review.Status() != ReviewCompleted {
			t.Errorf("expected status %s, got %s", ReviewCompleted, review.Status())
		}
		if review.CompletedAt() == nil {
			t.Error("expected completedAt to be set")
		}
	})

	t.Run("rejects completing an already completed review", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()
		review.RecordTaskCounts(1, 1)
		review.Complete()

		result := review.Complete()
		if result.IsOk() {
			t.Error("expected error for double completion")
		}
		if result.Error() != ErrAlreadyCompleted {
			t.Errorf("expected %q, got %q", ErrAlreadyCompleted, result.Error())
		}
	})
}

func TestDailyReview_Skip(t *testing.T) {
	t.Run("skips a pending review", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()

		result := review.Skip()
		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		if review.Status() != ReviewSkipped {
			t.Errorf("expected status %s, got %s", ReviewSkipped, review.Status())
		}
	})

	t.Run("rejects skipping a completed review", func(t *testing.T) {
		review := NewDailyReview(time.Now()).Value()
		review.RecordTaskCounts(1, 1)
		review.Complete()

		result := review.Skip()
		if result.IsOk() {
			t.Error("expected error for skipping completed review")
		}
	})
}
