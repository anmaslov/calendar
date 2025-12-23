package handler

import (
	"time"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/google/uuid"
)

// EventResponse represents an event in API response.
type EventResponse struct {
	ID          uuid.UUID  `json:"id"`
	ExchangeID  string     `json:"exchange_id"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body,omitempty"`
	Location    string     `json:"location,omitempty"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	IsAllDay    bool       `json:"is_all_day"`
	Organizer   string     `json:"organizer,omitempty"`
	Attendees   []string   `json:"attendees,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
	Importance  string     `json:"importance,omitempty"`
	Sensitivity string     `json:"sensitivity,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SyncedAt    *time.Time `json:"synced_at,omitempty"`
}

// ListEventsResponse represents the response for listing events.
type ListEventsResponse struct {
	Events []*EventResponse `json:"events"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
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

// toEventResponse converts domain event to API response.
func toEventResponse(e *domain.Event) *EventResponse {
	return &EventResponse{
		ID:          e.ID,
		ExchangeID:  e.ExchangeID,
		Subject:     e.Subject,
		Body:        e.Body,
		Location:    e.Location,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		IsAllDay:    e.IsAllDay,
		Organizer:   e.Organizer,
		Attendees:   e.Attendees,
		Categories:  e.Categories,
		Importance:  e.Importance,
		Sensitivity: e.Sensitivity,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		SyncedAt:    e.SyncedAt,
	}
}

// toEventResponseList converts a list of domain events to API responses.
func toEventResponseList(events []*domain.Event) []*EventResponse {
	result := make([]*EventResponse, len(events))
	for i, e := range events {
		result[i] = toEventResponse(e)
	}
	return result
}
