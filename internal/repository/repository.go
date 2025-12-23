package repository

import (
	"context"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/google/uuid"
)

// EventRepository defines the interface for event data access (read-only).
type EventRepository interface {
	// GetByID retrieves an event by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)

	// List retrieves events based on filter criteria.
	List(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error)

	// Count returns the total number of events matching the filter.
	Count(ctx context.Context, filter domain.EventFilter) (int64, error)
}

// EventSyncRepository defines the interface for event sync operations (write).
type EventSyncRepository interface {
	// Upsert creates or updates an event based on Exchange ID.
	Upsert(ctx context.Context, event *domain.Event) error

	// DeleteNotInExchangeIDs deletes events not in the provided Exchange IDs list.
	DeleteNotInExchangeIDs(ctx context.Context, exchangeIDs []string) error
}
