package habit

import (
	"time"

	domhabit "github.com/i-nishimura/goatodo/domain/habit"
	"github.com/i-nishimura/goatodo/domain/shared"
)

type GetStreak struct {
	repo domhabit.Repository
}

func NewGetStreak(repo domhabit.Repository) *GetStreak {
	return &GetStreak{repo: repo}
}

func (uc *GetStreak) Execute(today time.Time) shared.Result[domhabit.Streak] {
	// Fetch last 365 days of reviews for streak calculation
	from := today.AddDate(0, 0, -365)
	result := uc.repo.FindByDateRange(from, today)
	if result.IsErr() {
		return shared.Err[domhabit.Streak](result.Error())
	}

	streak := domhabit.CalculateStreak(result.Value(), today)
	return shared.Ok(streak)
}
