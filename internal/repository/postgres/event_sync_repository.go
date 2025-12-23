package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/anmaslov/calendar/internal/domain"
	"github.com/anmaslov/calendar/internal/repository"
	"github.com/jmoiron/sqlx"
)

// Event columns for sync operations
var eventColumns = []string{
	"id", "exchange_id", "subject", "body", "location",
	"start_time", "end_time", "is_all_day", "organizer",
	"importance", "sensitivity", "status", "created_at", "updated_at", "synced_at",
}

type eventSyncRepository struct {
	db *sqlx.DB
}

// NewEventSyncRepository creates a new PostgreSQL event sync repository.
func NewEventSyncRepository(db *sqlx.DB) repository.EventSyncRepository {
	return &eventSyncRepository{db: db}
}

func (r *eventSyncRepository) Upsert(ctx context.Context, event *domain.Event) error {
	now := time.Now()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	event.UpdatedAt = now
	event.SyncedAt = &now

	model := toEventModel(event)

	query, args, err := psql.Insert(eventsTable).
		Columns(eventColumns...).
		Values(
			model.ID, model.ExchangeID, model.Subject, model.Body, model.Location,
			model.StartTime, model.EndTime, model.IsAllDay, model.Organizer,
			model.Importance, model.Sensitivity, model.Status,
			model.CreatedAt, model.UpdatedAt, model.SyncedAt,
		).
		Suffix(`ON CONFLICT (exchange_id) DO UPDATE SET
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
			synced_at = EXCLUDED.synced_at`).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *eventSyncRepository) DeleteNotInExchangeIDs(ctx context.Context, exchangeIDs []string) error {
	builder := psql.Delete(eventsTable)
	if len(exchangeIDs) > 0 {
		builder = builder.Where(sq.NotEq{"exchange_id": exchangeIDs})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
