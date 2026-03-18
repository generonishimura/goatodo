package task

import (
	"time"

	"github.com/i-nishimura/goatodo/domain/shared"
)

type Task struct {
	id          string
	title       string
	description string
	status      Status
	priority    Priority
	createdAt   time.Time
	completedAt *time.Time
}

func NewTask(title string) shared.Result[*Task] {
	if title == "" {
		return shared.Err[*Task](ErrEmptyTitle)
	}
	return shared.Ok(&Task{
		id:        shared.NewID(),
		title:     title,
		status:    StatusTodo,
		priority:  PriorityNone,
		createdAt: time.Now(),
	})
}

func (t *Task) ID() string          { return t.id }
func (t *Task) Title() string       { return t.title }
func (t *Task) Description() string { return t.description }
func (t *Task) Status() Status      { return t.status }
func (t *Task) Priority() Priority  { return t.priority }
func (t *Task) CreatedAt() time.Time { return t.createdAt }
func (t *Task) CompletedAt() *time.Time { return t.completedAt }

var validTransitions = map[Status]map[Status]bool{
	StatusTodo:  {StatusTodo: true, StatusDoing: true, StatusDone: true},
	StatusDoing: {StatusDoing: true, StatusTodo: true, StatusDone: true},
	StatusDone:  {StatusDone: true},
}

func (t *Task) TransitionTo(newStatus Status) shared.Result[bool] {
	if !IsValidStatus(newStatus) {
		return shared.Err[bool](ErrInvalidStatus)
	}
	if allowed, ok := validTransitions[t.status][newStatus]; !ok || !allowed {
		return shared.Err[bool](ErrInvalidTransition)
	}
	t.status = newStatus
	if newStatus == StatusDone && t.completedAt == nil {
		now := time.Now()
		t.completedAt = &now
	}
	return shared.Ok(true)
}

func (t *Task) SetPriority(p Priority) shared.Result[bool] {
	if !IsValidPriority(p) {
		return shared.Err[bool](ErrInvalidPriority)
	}
	t.priority = p
	return shared.Ok(true)
}

func (t *Task) UpdateTitle(title string) shared.Result[bool] {
	if title == "" {
		return shared.Err[bool](ErrEmptyTitle)
	}
	t.title = title
	return shared.Ok(true)
}

func (t *Task) SetDescription(desc string) {
	t.description = desc
}

func Reconstruct(id, title, description string, status Status, priority Priority, createdAt time.Time, completedAt *time.Time) *Task {
	return &Task{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		createdAt:   createdAt,
		completedAt: completedAt,
	}
}
