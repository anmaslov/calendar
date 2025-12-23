package postgres

import (
	"time"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/google/uuid"
)

// eventModel represents a database model for event.
type eventModel struct {
	ID          uuid.UUID  `db:"id"`
	ExchangeID  string     `db:"exchange_id"`
	Subject     string     `db:"subject"`
	Body        string     `db:"body"`
	Location    string     `db:"location"`
	StartTime   time.Time  `db:"start_time"`
	EndTime     time.Time  `db:"end_time"`
	IsAllDay    bool       `db:"is_all_day"`
	Organizer   string     `db:"organizer"`
	Importance  string     `db:"importance"`
	Sensitivity string     `db:"sensitivity"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	SyncedAt    *time.Time `db:"synced_at"`
}

// toDomain converts database model to domain entity.
func (m *eventModel) toDomain() *domain.Event {
	return &domain.Event{
		ID:          m.ID,
		ExchangeID:  m.ExchangeID,
		Subject:     m.Subject,
		Body:        m.Body,
		Location:    m.Location,
		StartTime:   m.StartTime,
		EndTime:     m.EndTime,
		IsAllDay:    m.IsAllDay,
		Organizer:   m.Organizer,
		Importance:  m.Importance,
		Sensitivity: m.Sensitivity,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		SyncedAt:    m.SyncedAt,
	}
}

// toEventModel converts domain entity to database model.
func toEventModel(e *domain.Event) *eventModel {
	return &eventModel{
		ID:          e.ID,
		ExchangeID:  e.ExchangeID,
		Subject:     e.Subject,
		Body:        e.Body,
		Location:    e.Location,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		IsAllDay:    e.IsAllDay,
		Organizer:   e.Organizer,
		Importance:  e.Importance,
		Sensitivity: e.Sensitivity,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		SyncedAt:    e.SyncedAt,
	}
}
