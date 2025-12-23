package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
)

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	filter := parseEventFilter(r.URL.Query())

	events, total, err := h.eventService.ListEvents(r.Context(), filter)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list events")
		return
	}

	h.respondJSON(w, http.StatusOK, ListEventsResponse{
		Events: toEventResponseList(events),
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func parseEventFilter(q url.Values) domain.EventFilter {
	filter := domain.EventFilter{
		Limit:  defaultLimit,
		Offset: defaultOffset,
	}

	if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 {
		filter.Limit = min(l, maxLimit)
	}

	if o, err := strconv.Atoi(q.Get("offset")); err == nil && o >= 0 {
		filter.Offset = o
	}

	if t, err := time.Parse(time.RFC3339, q.Get("start_date")); err == nil {
		filter.StartDate = &t
	}

	if t, err := time.Parse(time.RFC3339, q.Get("end_date")); err == nil {
		filter.EndDate = &t
	}

	filter.Subject = q.Get("subject")
	filter.Status = q.Get("status")

	return filter
}

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.respondError(w, http.StatusBadRequest, "INVALID_ID", "Event ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID format")
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

	h.respondJSON(w, http.StatusOK, toEventResponse(event))
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
