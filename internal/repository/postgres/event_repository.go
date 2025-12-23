package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/anmaslov/calendar/internal/domain"
	"github.com/anmaslov/calendar/internal/repository"
	"github.com/jmoiron/sqlx"
)

// PostgreSQL placeholder format
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const eventsTable = "events"

type eventRepository struct {
	db *sqlx.DB
}

// NewEventRepository creates a new PostgreSQL event repository.
func NewEventRepository(db *sqlx.DB) repository.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	query, args, err := psql.Select("*").From(eventsTable).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	var event domain.Event
	if err := r.db.GetContext(ctx, &event, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) List(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	builder := applyEventFilter(psql.Select("*").From(eventsTable), filter).
		OrderBy("start_time ASC")

	if filter.Limit > 0 {
		builder = builder.Limit(uint64(filter.Limit))
	}

	if filter.Offset > 0 {
		builder = builder.Offset(uint64(filter.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var events []*domain.Event
	if err := r.db.SelectContext(ctx, &events, query, args...); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Count(ctx context.Context, filter domain.EventFilter) (int64, error) {
	query, args, err := applyEventFilter(psql.Select("COUNT(*)").From(eventsTable), filter).ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

// applyEventFilter applies common filters to the query builder.
func applyEventFilter(b sq.SelectBuilder, f domain.EventFilter) sq.SelectBuilder {
	if f.StartDate != nil {
		b = b.Where(sq.GtOrEq{"start_time": *f.StartDate})
	}
	if f.EndDate != nil {
		b = b.Where(sq.LtOrEq{"end_time": *f.EndDate})
	}
	if f.Subject != "" {
		b = b.Where(sq.ILike{"subject": "%" + f.Subject + "%"})
	}
	if f.Status != "" {
		b = b.Where(sq.Eq{"status": f.Status})
	}
	return b
}
