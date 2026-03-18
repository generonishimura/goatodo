package habit

import (
	"time"

	domhabit "github.com/i-nishimura/goatodo/domain/habit"
	"github.com/i-nishimura/goatodo/domain/shared"
	domtask "github.com/i-nishimura/goatodo/domain/task"
)

type CompleteDailyReview struct {
	habitRepo domhabit.Repository
	taskRepo  domtask.Repository
}

func NewCompleteDailyReview(habitRepo domhabit.Repository, taskRepo domtask.Repository) *CompleteDailyReview {
	return &CompleteDailyReview{habitRepo: habitRepo, taskRepo: taskRepo}
}

func (uc *CompleteDailyReview) Execute(date time.Time) shared.Result[*domhabit.DailyReview] {
	// Find or create review for the date
	findResult := uc.habitRepo.FindByDate(date)
	var review *domhabit.DailyReview
	if findResult.IsOk() {
		review = findResult.Value()
	} else {
		newResult := domhabit.NewDailyReview(date)
		if newResult.IsErr() {
			return newResult
		}
		review = newResult.Value()
	}

	// Count tasks
	tasksResult := uc.taskRepo.FindAll()
	completed, total := 0, 0
	if tasksResult.IsOk() {
		tasks := tasksResult.Value()
		total = len(tasks)
		for _, t := range tasks {
			if t.Status() == domtask.StatusDone {
				completed++
			}
		}
	}

	review.RecordTaskCounts(completed, total)

	// Complete the review
	completeResult := review.Complete()
	if completeResult.IsErr() {
		return shared.Err[*domhabit.DailyReview](completeResult.Error())
	}

	// Persist
	saveResult := uc.habitRepo.Save(review)
	if saveResult.IsErr() {
		return shared.Err[*domhabit.DailyReview](saveResult.Error())
	}

	return shared.Ok(review)
}
