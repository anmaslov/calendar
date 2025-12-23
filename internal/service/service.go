package service

import (
	"context"

	"github.com/anmaslov/calendar/internal/domain"
)

// EventService defines the interface for event business logic.
type EventService interface {
	// GetEvent retrieves an event by its ID.
	GetEvent(ctx context.Context, id string) (*domain.Event, error)

	// ListEvents retrieves events based on filter criteria.
	ListEvents(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, int64, error)
}
