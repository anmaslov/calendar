package sync

import (
	"context"
	"time"

	"github.com/anmaslov/calendar/internal/config"
	"github.com/anmaslov/calendar/internal/domain"
	"go.uber.org/zap"
)

// ExchangeClientStub is a stub implementation of ExchangeClient.
// Replace this with actual EWS implementation.
type ExchangeClientStub struct {
	cfg    config.ExchangeConfig
	logger *zap.Logger
}

// NewExchangeClient creates a new Exchange client.
// TODO: Replace with actual EWS implementation.
func NewExchangeClient(cfg config.ExchangeConfig, logger *zap.Logger) ExchangeClient {
	return &ExchangeClientStub{
		cfg:    cfg,
		logger: logger,
	}
}

func (c *ExchangeClientStub) GetCalendarEvents(ctx context.Context, startDate, endDate time.Time) ([]*domain.Event, error) {
	c.logger.Warn("using stub Exchange client - no real events will be fetched",
		zap.String("url", c.cfg.URL),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	// TODO: Implement actual EWS communication
	// Example using go-ews library or direct SOAP requests

	return []*domain.Event{}, nil
}

