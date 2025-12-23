package service

import (
	"context"

	"github.com/anmaslov/calendar/internal/domain"
)

// EventService defines the interface for event business logic.
type EventService interface {
	// CreateEvent creates a new event.
	CreateEvent(ctx context.Context, event *domain.Event) error

	// GetEvent retrieves an event by its ID.
	GetEvent(ctx context.Context, id string) (*domain.Event, error)

	// ListEvents retrieves events based on filter criteria.
	ListEvents(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, int64, error)

	// UpdateEvent updates an existing event.
	UpdateEvent(ctx context.Context, event *domain.Event) error

	// DeleteEvent deletes an event by its ID.
	DeleteEvent(ctx context.Context, id string) error

	// SyncEvents syncs events from Exchange server.
	SyncEvents(ctx context.Context) error
}

