package habit

import (
	"testing"
	"time"
)

func TestCalculateStreak(t *testing.T) {
	t.Run("returns zero streak when no reviews exist", func(t *testing.T) {
		streak := CalculateStreak(nil, time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC))

		if streak.Current != 0 {
			t.Errorf("expected current streak 0, got %d", streak.Current)
		}
		if streak.Longest != 0 {
			t.Errorf("expected longest streak 0, got %d", streak.Longest)
		}
	})

	t.Run("returns 1 when only today is completed", func(t *testing.T) {
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		reviews := []*DailyReview{
			makeCompletedReview(today),
		}

		streak := CalculateStreak(reviews, today)

		if streak.Current != 1 {
			t.Errorf("expected current streak 1, got %d", streak.Current)
		}
		if streak.Longest != 1 {
			t.Errorf("expected longest streak 1, got %d", streak.Longest)
		}
	})

	t.Run("counts consecutive completed days", func(t *testing.T) {
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		reviews := []*DailyReview{
			makeCompletedReview(today),
			makeCompletedReview(today.AddDate(0, 0, -1)),
			makeCompletedReview(today.AddDate(0, 0, -2)),
		}

		streak := CalculateStreak(reviews, today)

		if streak.Current != 3 {
			t.Errorf("expected current streak 3, got %d", streak.Current)
		}
	})

	t.Run("breaks streak on missed day", func(t *testing.T) {
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		reviews := []*DailyReview{
			makeCompletedReview(today),
			makeCompletedReview(today.AddDate(0, 0, -1)),
			// day -2 is missing
			makeCompletedReview(today.AddDate(0, 0, -3)),
			makeCompletedReview(today.AddDate(0, 0, -4)),
			makeCompletedReview(today.AddDate(0, 0, -5)),
		}

		streak := CalculateStreak(reviews, today)

		if streak.Current != 2 {
			t.Errorf("expected current streak 2, got %d", streak.Current)
		}
		if streak.Longest != 3 {
			t.Errorf("expected longest streak 3, got %d", streak.Longest)
		}
	})

	t.Run("skipped day breaks the streak", func(t *testing.T) {
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		skipped := NewDailyReview(today.AddDate(0, 0, -1)).Value()
		skipped.Skip()

		reviews := []*DailyReview{
			makeCompletedReview(today),
			skipped,
			makeCompletedReview(today.AddDate(0, 0, -2)),
		}

		streak := CalculateStreak(reviews, today)

		if streak.Current != 1 {
			t.Errorf("expected current streak 1, got %d", streak.Current)
		}
	})

	t.Run("current streak is 0 if today is not completed", func(t *testing.T) {
		today := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
		reviews := []*DailyReview{
			makeCompletedReview(today.AddDate(0, 0, -1)),
			makeCompletedReview(today.AddDate(0, 0, -2)),
		}

		streak := CalculateStreak(reviews, today)

		if streak.Current != 0 {
			t.Errorf("expected current streak 0, got %d", streak.Current)
		}
		if streak.Longest != 2 {
			t.Errorf("expected longest streak 2, got %d", streak.Longest)
		}
	})
}

func makeCompletedReview(date time.Time) *DailyReview {
	review := NewDailyReview(date).Value()
	review.RecordTaskCounts(1, 1)
	review.Complete()
	return review
}
