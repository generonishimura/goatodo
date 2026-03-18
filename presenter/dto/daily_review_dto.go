package dto

import (
	"github.com/i-nishimura/goatodo/domain/habit"
)

type DailyReviewDTO struct {
	ID                 string `json:"id"`
	Date               string `json:"date"`
	Status             string `json:"status"`
	CompletedTaskCount int    `json:"completedTaskCount"`
	TotalTaskCount     int    `json:"totalTaskCount"`
	CompletedAt        string `json:"completedAt,omitempty"`
}

type StreakDTO struct {
	Current int `json:"current"`
	Longest int `json:"longest"`
}

func FromDailyReview(r *habit.DailyReview) DailyReviewDTO {
	dto := DailyReviewDTO{
		ID:                 r.ID(),
		Date:               r.Date().Format("2006-01-02"),
		Status:             string(r.Status()),
		CompletedTaskCount: r.CompletedTaskCount(),
		TotalTaskCount:     r.TotalTaskCount(),
	}
	if r.CompletedAt() != nil {
		dto.CompletedAt = r.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
	}
	return dto
}

func FromStreak(s habit.Streak) StreakDTO {
	return StreakDTO{
		Current: s.Current,
		Longest: s.Longest,
	}
}
