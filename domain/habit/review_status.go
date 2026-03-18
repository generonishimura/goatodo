package habit

type ReviewStatus string

const (
	ReviewPending   ReviewStatus = "pending"
	ReviewCompleted ReviewStatus = "completed"
	ReviewSkipped   ReviewStatus = "skipped"
)

func IsValidReviewStatus(s ReviewStatus) bool {
	switch s {
	case ReviewPending, ReviewCompleted, ReviewSkipped:
		return true
	}
	return false
}
