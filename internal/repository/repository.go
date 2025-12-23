package repository

import (
	"context"

	"github.com/anmaslov/calendar/internal/domain"
)

// EventRepository defines the interface for event data access.
type EventRepository interface {
	// Create creates a new event.
	Create(ctx context.Context, event *domain.Event) error

	// GetByID retrieves an event by its ID.
	GetByID(ctx context.Context, id string) (*domain.Event, error)

	// GetByExchangeID retrieves an event by its Exchange ID.
	GetByExchangeID(ctx context.Context, exchangeID string) (*domain.Event, error)

	// List retrieves events based on filter criteria.
	List(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error)

	// Update updates an existing event.
	Update(ctx context.Context, event *domain.Event) error

	// Delete deletes an event by its ID.
	Delete(ctx context.Context, id string) error

	// Upsert creates or updates an event based on Exchange ID.
	Upsert(ctx context.Context, event *domain.Event) error

	// Count returns the total number of events matching the filter.
	Count(ctx context.Context, filter domain.EventFilter) (int64, error)
}

