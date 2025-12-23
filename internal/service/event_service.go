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

func (s *eventService) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get event", zap.String("id", id.String()), zap.Error(err))
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
