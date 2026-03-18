package habit

const (
	ErrInvalidTaskCount = "task count must be non-negative and completed must not exceed total"
	ErrAlreadyCompleted = "daily review is already completed"
	ErrNotPending       = "daily review is not in pending status"
)
