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
}

// New creates a new Handler.
func New(eventService service.EventService, logger *zap.Logger) *Handler {
	return &Handler{
		eventService: eventService,
		logger:       logger,
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

	// Health check
	r.Get("/health", h.healthCheck)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.listEvents)
			r.Post("/", h.createEvent)
			r.Get("/{id}", h.getEvent)
			r.Put("/{id}", h.updateEvent)
			r.Delete("/{id}", h.deleteEvent)
		})

		r.Post("/sync", h.syncEvents)
	})

	return r
}

