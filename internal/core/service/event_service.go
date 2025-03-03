package service

import (
	"context"
	"time"

	"astralis/internal/core/domain"
	"astralis/internal/core/ports"
)

type eventService struct {
	repositories []ports.EventRepository
}

// NewEventService creates a new instance of EventService
func NewEventService(repositories []ports.EventRepository) ports.EventService {
	return &eventService{
		repositories: repositories,
	}
}

// GetUpcomingEvents retrieves upcoming events from all repositories
func (s *eventService) GetUpcomingEvents(ctx context.Context, timeRange domain.TimeRange) ([]domain.Event, error) {
	var allEvents []domain.Event
	
	for _, repo := range s.repositories {
		events, err := repo.GetEvents(ctx, timeRange)
		if err != nil {
			continue // Skip failed repository but continue with others
		}
		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}

// GetEventByID retrieves a specific event by its ID from all repositories
func (s *eventService) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	for _, repo := range s.repositories {
		event, err := repo.GetEventByID(ctx, id)
		if err == nil && event != nil {
			return event, nil
		}
	}
	return nil, nil
}

// GetEventsByDate retrieves events for a specific date
func (s *eventService) GetEventsByDate(ctx context.Context, date time.Time) ([]domain.Event, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	return s.GetEventsByDateRange(ctx, startOfDay, endOfDay)
}

// GetEventsByDateRange retrieves events within a date range
func (s *eventService) GetEventsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error) {
	timeRange := domain.TimeRange{
		Start: start,
		End:   end,
	}
	return s.GetUpcomingEvents(ctx, timeRange)
}

// GetEventsByType retrieves events of a specific type
func (s *eventService) GetEventsByType(ctx context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
	var typedEvents []domain.Event
	
	for _, repo := range s.repositories {
		events, err := repo.GetEventsByType(ctx, eventType, timeRange)
		if err != nil {
			continue
		}
		typedEvents = append(typedEvents, events...)
	}

	return typedEvents, nil
}
