package presenter

import (
	"time"

	apphabit "github.com/i-nishimura/goatodo/application/habit"
	domhabit "github.com/i-nishimura/goatodo/domain/habit"
	domtask "github.com/i-nishimura/goatodo/domain/task"
	"github.com/i-nishimura/goatodo/presenter/dto"
)

type HabitHandler struct {
	completeDailyReview *apphabit.CompleteDailyReview
	getStreak           *apphabit.GetStreak
}

func NewHabitHandler(habitRepo domhabit.Repository, taskRepo domtask.Repository) *HabitHandler {
	return &HabitHandler{
		completeDailyReview: apphabit.NewCompleteDailyReview(habitRepo, taskRepo),
		getStreak:           apphabit.NewGetStreak(habitRepo),
	}
}

func (h *HabitHandler) CompleteDailyReview() TaskResponse {
	today := time.Now()
	result := h.completeDailyReview.Execute(today)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true, Data: dto.FromDailyReview(result.Value())}
}

func (h *HabitHandler) GetStreak() TaskResponse {
	today := time.Now()
	result := h.getStreak.Execute(today)
	if result.IsErr() {
		return TaskResponse{Success: false, Error: result.Error()}
	}
	return TaskResponse{Success: true, Data: dto.FromStreak(result.Value())}
}
