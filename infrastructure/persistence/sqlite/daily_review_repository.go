package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/i-nishimura/goatodo/domain/habit"
	"github.com/i-nishimura/goatodo/domain/shared"
)

type DailyReviewRepository struct {
	db *sql.DB
}

func NewDailyReviewRepository(db *sql.DB) *DailyReviewRepository {
	return &DailyReviewRepository{db: db}
}

func (r *DailyReviewRepository) Save(review *habit.DailyReview) shared.Result[bool] {
	var completedAt *string
	if review.CompletedAt() != nil {
		s := review.CompletedAt().Format(timeFormat)
		completedAt = &s
	}

	_, err := r.db.Exec(`
		INSERT INTO daily_reviews (id, date, status, completed_task_count, total_task_count, completed_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			status = excluded.status,
			completed_task_count = excluded.completed_task_count,
			total_task_count = excluded.total_task_count,
			completed_at = excluded.completed_at
	`,
		review.ID(),
		review.Date().Format(timeFormat),
		string(review.Status()),
		review.CompletedTaskCount(),
		review.TotalTaskCount(),
		completedAt,
	)
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	return shared.Ok(true)
}

func (r *DailyReviewRepository) FindByDate(date time.Time) shared.Result[*habit.DailyReview] {
	normalized := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	row := r.db.QueryRow(`
		SELECT id, date, status, completed_task_count, total_task_count, completed_at
		FROM daily_reviews WHERE date = ?
	`, normalized.Format(timeFormat))
	return scanReview(row)
}

func (r *DailyReviewRepository) FindByDateRange(from, to time.Time) shared.Result[[]*habit.DailyReview] {
	fromNorm := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toNorm := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)

	rows, err := r.db.Query(`
		SELECT id, date, status, completed_task_count, total_task_count, completed_at
		FROM daily_reviews
		WHERE date >= ? AND date <= ?
		ORDER BY date DESC
	`, fromNorm.Format(timeFormat), toNorm.Format(timeFormat))
	if err != nil {
		return shared.Err[[]*habit.DailyReview](err.Error())
	}
	defer rows.Close()

	var reviews []*habit.DailyReview
	for rows.Next() {
		result := scanReviewFromRows(rows)
		if result.IsErr() {
			return shared.Err[[]*habit.DailyReview](result.Error())
		}
		reviews = append(reviews, result.Value())
	}
	if err := rows.Err(); err != nil {
		return shared.Err[[]*habit.DailyReview](fmt.Sprintf("error iterating daily reviews: %v", err))
	}
	return shared.Ok(reviews)
}

func (r *DailyReviewRepository) Delete(id string) shared.Result[bool] {
	result, err := r.db.Exec("DELETE FROM daily_reviews WHERE id = ?", id)
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return shared.Err[bool](err.Error())
	}
	if affected == 0 {
		return shared.Err[bool]("daily review not found")
	}
	return shared.Ok(true)
}

func scanReview(row *sql.Row) shared.Result[*habit.DailyReview] {
	return scanReviewRow(row)
}

func scanReviewFromRows(rows *sql.Rows) shared.Result[*habit.DailyReview] {
	return scanReviewRow(rows)
}

func scanReviewRow(s scanner) shared.Result[*habit.DailyReview] {
	var (
		id, dateStr, status    string
		completedCount, totalCount int
		completedAtStr         *string
	)
	err := s.Scan(&id, &dateStr, &status, &completedCount, &totalCount, &completedAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.Err[*habit.DailyReview]("daily review not found")
		}
		return shared.Err[*habit.DailyReview](fmt.Sprintf("failed to scan daily review: %v", err))
	}

	date, err := time.Parse(timeFormat, dateStr)
	if err != nil {
		return shared.Err[*habit.DailyReview](fmt.Sprintf("failed to parse date: %v", err))
	}

	var completedAt *time.Time
	if completedAtStr != nil {
		t, err := time.Parse(timeFormat, *completedAtStr)
		if err != nil {
			return shared.Err[*habit.DailyReview](fmt.Sprintf("failed to parse completed_at: %v", err))
		}
		completedAt = &t
	}

	return shared.Ok(habit.Reconstruct(
		id, date, habit.ReviewStatus(status),
		completedCount, totalCount, completedAt,
	))
}
