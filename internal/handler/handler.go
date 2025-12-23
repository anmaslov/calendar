package handler

import (
	"github.com/anmaslov/calendar/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Handler holds all HTTP handlers.
type Handler struct {
	eventService service.EventService
	logger       *zap.Logger
	probes       *Probes
}

// New creates a new Handler.
func New(eventService service.EventService, logger *zap.Logger, probes *Probes) *Handler {
	return &Handler{
		eventService: eventService,
		logger:       logger,
		probes:       probes,
	}
}

// Router returns the HTTP router with all routes configured.
func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Health check (legacy)
	r.Get("/health", h.healthCheck)

	// Kubernetes probes
	r.Get("/healthz", h.livenessProbe) // Liveness probe
	r.Get("/readyz", h.readinessProbe) // Readiness probe

	// API v1 routes (read-only)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.listEvents)
			r.Get("/{id}", h.getEvent)
		})
	})

	return r
}
