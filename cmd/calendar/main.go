package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anmaslov/calendar/internal/config"
	"github.com/anmaslov/calendar/internal/handler"
	"github.com/anmaslov/calendar/internal/repository/postgres"
	"github.com/anmaslov/calendar/internal/service"
	"github.com/anmaslov/calendar/internal/sync"
	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Initialize Kubernetes probes
	probes := handler.NewProbes()

	// Connect to database
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	// Initialize repositories
	eventRepo := postgres.NewEventRepository(db)
	eventSyncRepo := postgres.NewEventSyncRepository(db)

	// Initialize services
	eventService := service.NewEventService(eventRepo, logger)

	// Initialize HTTP handler
	h := handler.New(eventService, logger, probes)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      h.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Create context for background workers
	ctx, cancel := context.WithCancel(context.Background())

	// Start sync worker if enabled
	var syncWorker *sync.Worker
	if cfg.Sync.Enabled {
		exchangeClient := sync.NewExchangeClient(cfg.Exchange, logger)
		syncWorker = sync.NewWorker(eventSyncRepo, exchangeClient, cfg.Sync, logger)
		syncWorker.Start(ctx)
	} else {
		logger.Info("sync worker is disabled")
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting server", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Mark application as ready after startup
	probes.SetReady(true)
	logger.Info("application is ready to receive traffic")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("received shutdown signal, starting graceful shutdown...")

	// Mark as not ready to stop receiving new traffic
	probes.SetReady(false)
	logger.Info("marked as not ready, waiting for in-flight requests...")

	// Stop sync worker
	if syncWorker != nil {
		syncWorker.Stop()
	}

	// Cancel context for all background workers
	cancel()

	// Give load balancer time to stop sending traffic
	time.Sleep(5 * time.Second)

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	// Shutdown HTTP server
	logger.Info("shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Close database connection
	logger.Info("closing database connection...")
	if err := db.Close(); err != nil {
		logger.Error("database close error", zap.Error(err))
	}

	logger.Info("graceful shutdown completed")
}
