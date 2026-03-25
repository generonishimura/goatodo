package memory

import (
	"sort"
	"sync"
	"time"

	"github.com/i-nishimura/goatodo/domain/habit"
	"github.com/i-nishimura/goatodo/domain/shared"
)

type DailyReviewRepository struct {
	mu      sync.RWMutex
	reviews map[string]*habit.DailyReview
}

func NewDailyReviewRepository() *DailyReviewRepository {
	return &DailyReviewRepository{
		reviews: make(map[string]*habit.DailyReview),
	}
}

func dateKey(date time.Time) string {
	return date.Format("2006-01-02")
}

func (r *DailyReviewRepository) Save(review *habit.DailyReview) shared.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reviews[dateKey(review.Date())] = review
	return shared.Ok(true)
}

func (r *DailyReviewRepository) FindByDate(date time.Time) shared.Result[*habit.DailyReview] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	if review, ok := r.reviews[dateKey(normalized)]; ok {
		return shared.Ok(review)
	}
	return shared.Err[*habit.DailyReview]("daily review not found")
}

func (r *DailyReviewRepository) FindByDateRange(from, to time.Time) shared.Result[[]*habit.DailyReview] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fromNorm := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toNorm := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	result := make([]*habit.DailyReview, 0)
	for _, review := range r.reviews {
		d := review.Date()
		if (d.Equal(fromNorm) || d.After(fromNorm)) && (d.Equal(toNorm) || d.Before(toNorm)) {
			result = append(result, review)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date().After(result[j].Date())
	})
	return shared.Ok(result)
}

func (r *DailyReviewRepository) Delete(id string) shared.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, review := range r.reviews {
		if review.ID() == id {
			delete(r.reviews, key)
			return shared.Ok(true)
		}
	}
	return shared.Err[bool]("daily review not found")
}
