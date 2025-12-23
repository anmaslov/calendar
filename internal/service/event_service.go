package service

import (
	"context"

	"github.com/anmaslov/calendar/internal/domain"
	"github.com/anmaslov/calendar/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type eventService struct {
	repo   repository.EventRepository
	logger *zap.Logger
}

// NewEventService creates a new event service.
func NewEventService(repo repository.EventRepository, logger *zap.Logger) EventService {
	return &eventService{
		repo:   repo,
		logger: logger,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *domain.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	if err := s.validateEvent(event); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, event); err != nil {
		s.logger.Error("failed to create event", zap.Error(err))
		return err
	}

	s.logger.Info("event created", zap.String("id", event.ID))
	return nil
}

func (s *eventService) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get event", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return event, nil
}

func (s *eventService) ListEvents(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, int64, error) {
	events, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list events", zap.Error(err))
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		s.logger.Error("failed to count events", zap.Error(err))
		return nil, 0, err
	}

	return events, count, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, event *domain.Event) error {
	if err := s.validateEvent(event); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, event); err != nil {
		s.logger.Error("failed to update event", zap.String("id", event.ID), zap.Error(err))
		return err
	}

	s.logger.Info("event updated", zap.String("id", event.ID))
	return nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete event", zap.String("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("event deleted", zap.String("id", id))
	return nil
}

func (s *eventService) SyncEvents(ctx context.Context) error {
	// TODO: Implement Exchange sync logic
	s.logger.Info("sync events called - not implemented yet")
	return nil
}

func (s *eventService) validateEvent(event *domain.Event) error {
	if event.Subject == "" {
		return domain.ErrInvalidInput
	}

	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		return domain.ErrInvalidInput
	}

	if event.EndTime.Before(event.StartTime) {
		return domain.ErrInvalidInput
	}

	return nil
}

