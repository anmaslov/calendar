package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/go-chi/chi/v5"
)

// ListEventsResponse represents the response for listing events.
type ListEventsResponse struct {
	Events []*domain.Event `json:"events"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents error details.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	filter := domain.EventFilter{
		Limit:  20,
		Offset: 0,
	}

	// Parse query parameters
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filter.StartDate = &t
		}
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filter.EndDate = &t
		}
	}

	if subject := r.URL.Query().Get("subject"); subject != "" {
		filter.Subject = subject
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = status
	}

	events, total, err := h.eventService.ListEvents(r.Context(), filter)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list events")
		return
	}

	h.respondJSON(w, http.StatusOK, ListEventsResponse{
		Events: events,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "INVALID_ID", "Event ID is required")
		return
	}

	event, err := h.eventService.GetEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			h.respondError(w, http.StatusNotFound, "NOT_FOUND", "Event not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get event")
		return
	}

	h.respondJSON(w, http.StatusOK, event)
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response")
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, code, message string) {
	h.respondJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
