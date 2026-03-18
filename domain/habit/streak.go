package habit

import (
	"sort"
	"time"
)

type Streak struct {
	Current int
	Longest int
}

func CalculateStreak(reviews []*DailyReview, today time.Time) Streak {
	if len(reviews) == 0 {
		return Streak{}
	}

	completedDates := make(map[time.Time]bool)
	for _, r := range reviews {
		if r.Status() == ReviewCompleted {
			completedDates[r.Date()] = true
		}
	}

	if len(completedDates) == 0 {
		return Streak{}
	}

	todayNorm := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate current streak (must include today)
	current := 0
	if completedDates[todayNorm] {
		current = 1
		for d := todayNorm.AddDate(0, 0, -1); completedDates[d]; d = d.AddDate(0, 0, -1) {
			current++
		}
	}

	// Calculate longest streak
	dates := make([]time.Time, 0, len(completedDates))
	for d := range completedDates {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

	longest := 1
	streak := 1
	for i := 1; i < len(dates); i++ {
		if dates[i].Sub(dates[i-1]) == 24*time.Hour {
			streak++
			if streak > longest {
				longest = streak
			}
		} else {
			streak = 1
		}
	}

	if current > longest {
		longest = current
	}

	return Streak{Current: current, Longest: longest}
}
