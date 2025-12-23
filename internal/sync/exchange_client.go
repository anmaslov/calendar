package sync

import (
	"context"
	"time"

	"github.com/anmaslov/calendar/internal/domain"
)

// ExchangeClient defines the interface for Exchange server communication.
type ExchangeClient interface {
	// GetCalendarEvents fetches calendar events from Exchange server.
	GetCalendarEvents(ctx context.Context, startDate, endDate time.Time) ([]*domain.Event, error)
}

