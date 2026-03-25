package habit

import (
	"time"

	"github.com/i-nishimura/goatodo/domain/shared"
)

type DailyReview struct {
	id                 string
	date               time.Time
	status             ReviewStatus
	completedTaskCount int
	totalTaskCount     int
	completedAt        *time.Time
}

func NewDailyReview(date time.Time) shared.Result[*DailyReview] {
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return shared.Ok(&DailyReview{
		id:     shared.NewID(),
		date:   normalized,
		status: ReviewPending,
	})
}

func (r *DailyReview) ID() string                { return r.id }
func (r *DailyReview) Date() time.Time           { return r.date }
func (r *DailyReview) Status() ReviewStatus       { return r.status }
func (r *DailyReview) CompletedTaskCount() int    { return r.completedTaskCount }
func (r *DailyReview) TotalTaskCount() int        { return r.totalTaskCount }
func (r *DailyReview) CompletedAt() *time.Time    { return r.completedAt }

func (r *DailyReview) RecordTaskCounts(completed, total int) shared.Result[bool] {
	if completed < 0 || total < 0 || completed > total {
		return shared.Err[bool](ErrInvalidTaskCount)
	}
	r.completedTaskCount = completed
	r.totalTaskCount = total
	return shared.Ok(true)
}

func (r *DailyReview) Complete() shared.Result[bool] {
	if r.status == ReviewCompleted {
		return shared.Err[bool](ErrAlreadyCompleted)
	}
	if r.status != ReviewPending {
		return shared.Err[bool](ErrNotPending)
	}
	r.status = ReviewCompleted
	now := time.Now()
	r.completedAt = &now
	return shared.Ok(true)
}

func (r *DailyReview) Skip() shared.Result[bool] {
	if r.status != ReviewPending {
		return shared.Err[bool](ErrNotPending)
	}
	r.status = ReviewSkipped
	return shared.Ok(true)
}
