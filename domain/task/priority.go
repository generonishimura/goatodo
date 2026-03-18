package task

type Priority int

const (
	PriorityNone   Priority = 0
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

func IsValidPriority(p Priority) bool {
	return p >= PriorityNone && p <= PriorityHigh
}
