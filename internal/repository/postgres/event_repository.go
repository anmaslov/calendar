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
	query, args, err := psql.
		Select("*").
		From(eventsTable).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var event domain.Event
	err = r.db.GetContext(ctx, &event, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) List(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	builder := psql.
		Select("*").
		From(eventsTable)

	builder = applyEventFilter(builder, filter)
	builder = builder.OrderBy("start_time ASC")

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
	err = r.db.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Count(ctx context.Context, filter domain.EventFilter) (int64, error) {
	builder := psql.
		Select("COUNT(*)").
		From(eventsTable)

	builder = applyEventFilter(builder, filter)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// applyEventFilter applies common filters to the query builder.
func applyEventFilter(builder sq.SelectBuilder, filter domain.EventFilter) sq.SelectBuilder {
	if filter.StartDate != nil {
		builder = builder.Where(sq.GtOrEq{"start_time": *filter.StartDate})
	}

	if filter.EndDate != nil {
		builder = builder.Where(sq.LtOrEq{"end_time": *filter.EndDate})
	}

	if filter.Subject != "" {
		builder = builder.Where(sq.ILike{"subject": "%" + filter.Subject + "%"})
	}

	if filter.Status != "" {
		builder = builder.Where(sq.Eq{"status": filter.Status})
	}

	return builder
}
