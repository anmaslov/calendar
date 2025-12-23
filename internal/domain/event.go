package domain

import "time"

// Event represents a calendar event.
type Event struct {
	ID          string     `json:"id" db:"id"`
	ExchangeID  string     `json:"exchange_id" db:"exchange_id"`
	Subject     string     `json:"subject" db:"subject"`
	Body        string     `json:"body,omitempty" db:"body"`
	Location    string     `json:"location,omitempty" db:"location"`
	StartTime   time.Time  `json:"start_time" db:"start_time"`
	EndTime     time.Time  `json:"end_time" db:"end_time"`
	IsAllDay    bool       `json:"is_all_day" db:"is_all_day"`
	Organizer   string     `json:"organizer,omitempty" db:"organizer"`
	Attendees   []string   `json:"attendees,omitempty" db:"-"`
	Categories  []string   `json:"categories,omitempty" db:"-"`
	Importance  string     `json:"importance,omitempty" db:"importance"`
	Sensitivity string     `json:"sensitivity,omitempty" db:"sensitivity"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	SyncedAt    *time.Time `json:"synced_at,omitempty" db:"synced_at"`
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

