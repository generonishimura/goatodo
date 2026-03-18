package habit

import "time"

func Reconstruct(id string, date time.Time, status ReviewStatus, completedTaskCount, totalTaskCount int, completedAt *time.Time) *DailyReview {
	return &DailyReview{
		id:                 id,
		date:               date,
		status:             status,
		completedTaskCount: completedTaskCount,
		totalTaskCount:     totalTaskCount,
		completedAt:        completedAt,
	}
}
