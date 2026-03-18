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

func (r *DailyReviewRepository) Save(review *habit.DailyReview) shared.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reviews[review.ID()] = review
	return shared.Ok(true)
}

func (r *DailyReviewRepository) FindByDate(date time.Time) shared.Result[*habit.DailyReview] {
	r.mu.RLock()
	defer r.mu.RUnlock()
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	for _, review := range r.reviews {
		if review.Date().Equal(normalized) {
			return shared.Ok(review)
		}
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
	if _, ok := r.reviews[id]; !ok {
		return shared.Err[bool]("daily review not found")
	}
	delete(r.reviews, id)
	return shared.Ok(true)
}
