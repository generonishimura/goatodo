package habit

import (
	"testing"
	"time"

	domhabit "github.com/i-nishimura/goatodo/domain/habit"
	"github.com/i-nishimura/goatodo/infrastructure/persistence/memory"
)

func TestGetStreak_Execute(t *testing.T) {
	t.Run("returns zero streak when no reviews", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		uc := NewGetStreak(repo)

		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		result := uc.Execute(today)

		if result.IsErr() {
			t.Fatalf("expected ok, got error: %s", result.Error())
		}

		streak := result.Value()
		if streak.Current != 0 {
			t.Errorf("expected current streak 0, got %d", streak.Current)
		}
	})

	t.Run("returns streak from completed reviews", func(t *testing.T) {
		repo := memory.NewDailyReviewRepository()
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)

		for i := 0; i < 3; i++ {
			review := domhabit.NewDailyReview(today.AddDate(0, 0, -i)).Value()
			review.RecordTaskCounts(1, 1)
			review.Complete()
			repo.Save(review)
		}

		uc := NewGetStreak(repo)
		result := uc.Execute(today)

		streak := result.Value()
		if streak.Current != 3 {
			t.Errorf("expected current streak 3, got %d", streak.Current)
		}
		if streak.Longest != 3 {
			t.Errorf("expected longest streak 3, got %d", streak.Longest)
		}
	})
}
