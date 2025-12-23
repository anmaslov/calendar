package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/anmaslov/calendar/internal/repository"
	"github.com/jmoiron/sqlx"
)

type eventRepository struct {
	db *sqlx.DB
}

// NewEventRepository creates a new PostgreSQL event repository.
func NewEventRepository(db *sqlx.DB) repository.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (
			id, exchange_id, subject, body, location,
			start_time, end_time, is_all_day, organizer,
			importance, sensitivity, status, created_at, updated_at, synced_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)`

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.ExchangeID, event.Subject, event.Body, event.Location,
		event.StartTime, event.EndTime, event.IsAllDay, event.Organizer,
		event.Importance, event.Sensitivity, event.Status,
		event.CreatedAt, event.UpdatedAt, event.SyncedAt,
	)

	return err
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	query := `SELECT * FROM events WHERE id = $1`

	var event domain.Event
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) GetByExchangeID(ctx context.Context, exchangeID string) (*domain.Event, error) {
	query := `SELECT * FROM events WHERE exchange_id = $1`

	var event domain.Event
	err := r.db.GetContext(ctx, &event, query, exchangeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) List(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	query := `SELECT * FROM events WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND end_time <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.Subject != "" {
		query += fmt.Sprintf(" AND subject ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Subject+"%")
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	query += " ORDER BY start_time ASC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	var events []*domain.Event
	err := r.db.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events SET
			subject = $2,
			body = $3,
			location = $4,
			start_time = $5,
			end_time = $6,
			is_all_day = $7,
			organizer = $8,
			importance = $9,
			sensitivity = $10,
			status = $11,
			updated_at = $12,
			synced_at = $13
		WHERE id = $1`

	event.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		event.ID, event.Subject, event.Body, event.Location,
		event.StartTime, event.EndTime, event.IsAllDay, event.Organizer,
		event.Importance, event.Sensitivity, event.Status,
		event.UpdatedAt, event.SyncedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) Upsert(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (
			id, exchange_id, subject, body, location,
			start_time, end_time, is_all_day, organizer,
			importance, sensitivity, status, created_at, updated_at, synced_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (exchange_id) DO UPDATE SET
			subject = EXCLUDED.subject,
			body = EXCLUDED.body,
			location = EXCLUDED.location,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			is_all_day = EXCLUDED.is_all_day,
			organizer = EXCLUDED.organizer,
			importance = EXCLUDED.importance,
			sensitivity = EXCLUDED.sensitivity,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at,
			synced_at = EXCLUDED.synced_at`

	now := time.Now()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	event.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.ExchangeID, event.Subject, event.Body, event.Location,
		event.StartTime, event.EndTime, event.IsAllDay, event.Organizer,
		event.Importance, event.Sensitivity, event.Status,
		event.CreatedAt, event.UpdatedAt, event.SyncedAt,
	)

	return err
}

func (r *eventRepository) Count(ctx context.Context, filter domain.EventFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM events WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND end_time <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.Subject != "" {
		query += fmt.Sprintf(" AND subject ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Subject+"%")
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
	}

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

