package habit

import (
	"time"

	"github.com/i-nishimura/goatodo/domain/shared"
)

type Repository interface {
	Save(review *DailyReview) shared.Result[bool]
	FindByDate(date time.Time) shared.Result[*DailyReview]
	FindByDateRange(from, to time.Time) shared.Result[[]*DailyReview]
	Delete(id string) shared.Result[bool]
}
