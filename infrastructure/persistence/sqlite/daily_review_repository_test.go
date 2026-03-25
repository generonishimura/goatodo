package sqlite_test

import (
	"testing"
	"time"

	"github.com/i-nishimura/goatodo/domain/habit"
	sqliterepo "github.com/i-nishimura/goatodo/infrastructure/persistence/sqlite"
)

func TestSQLiteDailyReviewRepository_SaveAndFindByDate(t *testing.T) {
	t.Run("saves and retrieves a daily review by date", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewDailyReviewRepository(db)

		date := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		review := habit.NewDailyReview(date).Value()
		review.RecordTaskCounts(3, 5)
		review.Complete()

		saveResult := repo.Save(review)
		if saveResult.IsErr() {
			t.Fatalf("save failed: %s", saveResult.Error())
		}

		found := repo.FindByDate(date)
		if found.IsErr() {
			t.Fatalf("find failed: %s", found.Error())
		}

		got := found.Value()
		if got.ID() != review.ID() {
			t.Errorf("expected ID %s, got %s", review.ID(), got.ID())
		}
		if got.Status() != habit.ReviewCompleted {
			t.Errorf("expected status %s, got %s", habit.ReviewCompleted, got.Status())
		}
		if got.CompletedTaskCount() != 3 {
			t.Errorf("expected 3 completed tasks, got %d", got.CompletedTaskCount())
		}
		if got.TotalTaskCount() != 5 {
			t.Errorf("expected 5 total tasks, got %d", got.TotalTaskCount())
		}
		if got.CompletedAt() == nil {
			t.Error("expected completedAt to be set")
		}
	})
}

func TestSQLiteDailyReviewRepository_UpsertByDate(t *testing.T) {
	t.Run("upserts by date, updating status and counts", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewDailyReviewRepository(db)

		date := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)

		// Save a pending review
		review1 := habit.NewDailyReview(date).Value()
		repo.Save(review1)

		// Save another review for the same date (different ID) — should upsert
		review2 := habit.NewDailyReview(date).Value()
		review2.RecordTaskCounts(2, 4)
		review2.Complete()
		repo.Save(review2)

		// Should find the updated version
		found := repo.FindByDate(date)
		if found.IsErr() {
			t.Fatalf("find failed: %s", found.Error())
		}

		got := found.Value()
		if got.Status() != habit.ReviewCompleted {
			t.Errorf("expected status %s, got %s", habit.ReviewCompleted, got.Status())
		}
		if got.CompletedTaskCount() != 2 {
			t.Errorf("expected 2 completed tasks, got %d", got.CompletedTaskCount())
		}
	})
}

func TestSQLiteDailyReviewRepository_FindByDateRange(t *testing.T) {
	t.Run("returns reviews within date range", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewDailyReviewRepository(db)

		for i := 15; i <= 18; i++ {
			date := time.Date(2026, 3, i, 0, 0, 0, 0, time.UTC)
			review := habit.NewDailyReview(date).Value()
			review.RecordTaskCounts(1, 1)
			review.Complete()
			repo.Save(review)
		}

		from := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
		to := time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC)
		result := repo.FindByDateRange(from, to)
		if result.IsErr() {
			t.Fatalf("find range failed: %s", result.Error())
		}
		if len(result.Value()) != 2 {
			t.Errorf("expected 2 reviews, got %d", len(result.Value()))
		}
	})
}

func TestSQLiteDailyReviewRepository_FindByDateNotFound(t *testing.T) {
	t.Run("returns error for non-existent date", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewDailyReviewRepository(db)

		date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		result := repo.FindByDate(date)
		if result.IsOk() {
			t.Error("expected error for non-existent date")
		}
		if result.Error() != habit.ErrNotFound {
			t.Errorf("expected %q, got %q", habit.ErrNotFound, result.Error())
		}
	})
}

func TestSQLiteDailyReviewRepository_Delete(t *testing.T) {
	t.Run("deletes a daily review", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewDailyReviewRepository(db)

		date := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		review := habit.NewDailyReview(date).Value()
		repo.Save(review)

		delResult := repo.Delete(review.ID())
		if delResult.IsErr() {
			t.Fatalf("delete failed: %s", delResult.Error())
		}

		found := repo.FindByDate(date)
		if found.IsOk() {
			t.Error("review should have been deleted")
		}
	})
}
