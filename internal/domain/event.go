package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a calendar event.
type Event struct {
	ID          uuid.UUID
	ExchangeID  string
	Subject     string
	Body        string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	IsAllDay    bool
	Organizer   string
	Attendees   []string
	Categories  []string
	Importance  string
	Sensitivity string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SyncedAt    *time.Time
}

// NewEvent creates a new event with generated UUID.
func NewEvent() *Event {
	return &Event{
		ID: uuid.New(),
	}
}

// EventFilter represents filters for querying events.
type EventFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Subject   string
	Status    string
	Limit     int
	Offset    int
}
