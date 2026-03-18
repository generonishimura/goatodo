package sqlite

import (
	"database/sql"
	"time"

	"github.com/i-nishimura/goatodo/domain/shared"
	"github.com/i-nishimura/goatodo/domain/task"
)

const timeFormat = time.RFC3339

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(t *task.Task) shared.Result[bool] {
	var completedAt *string
	if t.CompletedAt() != nil {
		s := t.CompletedAt().Format(timeFormat)
		completedAt = &s
	}

	_, err := r.db.Exec(`
		INSERT INTO tasks (id, title, description, status, priority, created_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			status = excluded.status,
			priority = excluded.priority,
			completed_at = excluded.completed_at
	`,
		t.ID(),
		t.Title(),
		t.Description(),
		string(t.Status()),
		int(t.Priority()),
		t.CreatedAt().Format(timeFormat),
		completedAt,
	)
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	return shared.Ok(true)
}

func (r *TaskRepository) FindByID(id string) shared.Result[*task.Task] {
	row := r.db.QueryRow(`
		SELECT id, title, description, status, priority, created_at, completed_at
		FROM tasks WHERE id = ?
	`, id)
	return scanTask(row)
}

func (r *TaskRepository) FindAll() shared.Result[[]*task.Task] {
	rows, err := r.db.Query(`
		SELECT id, title, description, status, priority, created_at, completed_at
		FROM tasks ORDER BY created_at DESC
	`)
	if err != nil {
		return shared.Err[[]*task.Task](err.Error())
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		result := scanTaskFromRows(rows)
		if result.IsErr() {
			return shared.Err[[]*task.Task](result.Error())
		}
		tasks = append(tasks, result.Value())
	}
	return shared.Ok(tasks)
}

func (r *TaskRepository) Delete(id string) shared.Result[bool] {
	result, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	if affected == 0 {
		return shared.Err[bool]("task not found")
	}
	return shared.Ok(true)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRow(s scanner) shared.Result[*task.Task] {
	var (
		id, title, description, status, createdAtStr string
		priority                                     int
		completedAtStr                               *string
	)
	err := s.Scan(&id, &title, &description, &status, &priority, &createdAtStr, &completedAtStr)
	if err != nil {
		return shared.Err[*task.Task]("task not found")
	}

	createdAt, _ := time.Parse(timeFormat, createdAtStr)
	var completedAt *time.Time
	if completedAtStr != nil {
		t, _ := time.Parse(timeFormat, *completedAtStr)
		completedAt = &t
	}

	return shared.Ok(task.Reconstruct(
		id, title, description,
		task.Status(status),
		task.Priority(priority),
		createdAt, completedAt,
	))
}

func scanTask(row *sql.Row) shared.Result[*task.Task] {
	return scanRow(row)
}

func scanTaskFromRows(rows *sql.Rows) shared.Result[*task.Task] {
	return scanRow(rows)
}
