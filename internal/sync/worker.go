package sync

import (
	"context"
	"time"

	"github.com/anmaslov/calendar/internal/config"
	"github.com/anmaslov/calendar/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Worker handles background synchronization with Exchange server.
type Worker struct {
	syncRepo       repository.EventSyncRepository
	exchangeClient ExchangeClient
	cfg            config.SyncConfig
	logger         *zap.Logger
	stopCh         chan struct{}
	doneCh         chan struct{}
}

// NewWorker creates a new sync worker.
func NewWorker(
	syncRepo repository.EventSyncRepository,
	exchangeClient ExchangeClient,
	cfg config.SyncConfig,
	logger *zap.Logger,
) *Worker {
	return &Worker{
		syncRepo:       syncRepo,
		exchangeClient: exchangeClient,
		cfg:            cfg,
		logger:         logger,
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}
}

// Start starts the background sync worker.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("starting sync worker",
		zap.Duration("interval", w.cfg.Interval),
		zap.Int("sync_days", w.cfg.SyncDays),
	)

	go w.run(ctx)
}

// Stop stops the sync worker gracefully.
func (w *Worker) Stop() {
	w.logger.Info("stopping sync worker...")
	close(w.stopCh)
	<-w.doneCh
	w.logger.Info("sync worker stopped")
}

func (w *Worker) run(ctx context.Context) {
	defer close(w.doneCh)

	// Run initial sync
	w.sync(ctx)

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.sync(ctx)
		}
	}
}

func (w *Worker) sync(ctx context.Context) {
	w.logger.Info("starting sync cycle")
	startTime := time.Now()

	// Calculate date range for sync
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, w.cfg.SyncDays)

	// Fetch events from Exchange
	events, err := w.exchangeClient.GetCalendarEvents(ctx, startDate, endDate)
	if err != nil {
		w.logger.Error("failed to fetch events from Exchange", zap.Error(err))
		return
	}

	w.logger.Info("fetched events from Exchange", zap.Int("count", len(events)))

	// Collect Exchange IDs for cleanup
	exchangeIDs := make([]string, 0, len(events))

	// Upsert events
	for _, event := range events {
		// Generate UUID if not set
		if event.ID == uuid.Nil {
			event.ID = uuid.New()
		}

		if err := w.syncRepo.Upsert(ctx, event); err != nil {
			w.logger.Error("failed to upsert event",
				zap.String("exchange_id", event.ExchangeID),
				zap.Error(err),
			)
			continue
		}
		exchangeIDs = append(exchangeIDs, event.ExchangeID)
	}

	// Delete events that no longer exist in Exchange
	if err := w.syncRepo.DeleteNotInExchangeIDs(ctx, exchangeIDs); err != nil {
		w.logger.Error("failed to delete old events", zap.Error(err))
	}

	w.logger.Info("sync cycle completed",
		zap.Duration("duration", time.Since(startTime)),
		zap.Int("synced_events", len(exchangeIDs)),
	)
}
