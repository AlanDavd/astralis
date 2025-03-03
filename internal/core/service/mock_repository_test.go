package service

import (
	"context"
	"time"

	"astralis/internal/core/domain"
)

type mockRepository struct {
	events map[string]domain.Event
}

func newMockRepository() *mockRepository {
	now := time.Now()
	events := map[string]domain.Event{
		"meteor-1": {
			ID:          "meteor-1",
			Title:       "Test Meteor Shower",
			Description: "Test Description",
			StartTime:   now,
			EndTime:     now.Add(24 * time.Hour),
			Type:        domain.MeteorShower,
		},
		"eclipse-1": {
			ID:          "eclipse-1",
			Title:       "Test Eclipse",
			Description: "Test Description",
			StartTime:   now.Add(48 * time.Hour),
			EndTime:     now.Add(72 * time.Hour),
			Type:        domain.Eclipse,
		},
	}

	return &mockRepository{
		events: events,
	}
}

func (r *mockRepository) GetEvents(_ context.Context, timeRange domain.TimeRange) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range r.events {
		if event.StartTime.After(timeRange.Start) && event.StartTime.Before(timeRange.End) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (r *mockRepository) GetEventByID(_ context.Context, id string) (*domain.Event, error) {
	if event, ok := r.events[id]; ok {
		return &event, nil
	}
	return nil, nil
}

func (r *mockRepository) GetEventsByType(_ context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range r.events {
		if event.Type == eventType && event.StartTime.After(timeRange.Start) && event.StartTime.Before(timeRange.End) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (r *mockRepository) Name() string {
	return "Mock Repository"
} 