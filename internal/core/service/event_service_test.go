package service

import (
	"context"
	"testing"
	"time"

	"astralis/internal/core/domain"
	"astralis/internal/core/ports"
)

func TestEventService_GetUpcomingEvents(t *testing.T) {
	mockRepo := newMockRepository()
	service := NewEventService([]ports.EventRepository{mockRepo})

	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-1 * time.Hour),
		End:   now.Add(96 * time.Hour),
	}

	events, err := service.GetUpcomingEvents(context.Background(), timeRange)
	if err != nil {
		t.Errorf("GetUpcomingEvents() error = %v", err)
		return
	}

	if len(events) != 2 {
		t.Errorf("GetUpcomingEvents() got %v events, want %v", len(events), 2)
	}
}

func TestEventService_GetEventByID(t *testing.T) {
	mockRepo := newMockRepository()
	service := NewEventService([]ports.EventRepository{mockRepo})

	tests := []struct {
		name    string
		id      string
		wantID  string
		wantNil bool
	}{
		{
			name:    "existing event",
			id:      "meteor-1",
			wantID:  "meteor-1",
			wantNil: false,
		},
		{
			name:    "non-existing event",
			id:      "non-existing",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := service.GetEventByID(context.Background(), tt.id)
			if err != nil {
				t.Errorf("GetEventByID() error = %v", err)
				return
			}

			if tt.wantNil && event != nil {
				t.Errorf("GetEventByID() got event = %v, want nil", event)
				return
			}

			if !tt.wantNil && event == nil {
				t.Error("GetEventByID() got nil, want event")
				return
			}

			if !tt.wantNil && event.ID != tt.wantID {
				t.Errorf("GetEventByID() got event ID = %v, want %v", event.ID, tt.wantID)
			}
		})
	}
}

func TestEventService_GetEventsByType(t *testing.T) {
	mockRepo := newMockRepository()
	service := NewEventService([]ports.EventRepository{mockRepo})

	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-1 * time.Hour),
		End:   now.Add(96 * time.Hour),
	}

	tests := []struct {
		name      string
		eventType domain.EventType
		wantCount int
	}{
		{
			name:      "meteor shower events",
			eventType: domain.MeteorShower,
			wantCount: 1,
		},
		{
			name:      "eclipse events",
			eventType: domain.Eclipse,
			wantCount: 1,
		},
		{
			name:      "no events of type",
			eventType: domain.Transit,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsByType(context.Background(), tt.eventType, timeRange)
			if err != nil {
				t.Errorf("GetEventsByType() error = %v", err)
				return
			}

			if len(events) != tt.wantCount {
				t.Errorf("GetEventsByType() got %v events, want %v", len(events), tt.wantCount)
			}

			for _, event := range events {
				if event.Type != tt.eventType {
					t.Errorf("GetEventsByType() got event of type %v, want %v", event.Type, tt.eventType)
				}
			}
		})
	}
}

func TestEventService_GetEventsByDate(t *testing.T) {
	mockRepo := newMockRepository()
	service := NewEventService([]ports.EventRepository{mockRepo})

	now := time.Now()
	tests := []struct {
		name      string
		date      time.Time
		wantCount int
	}{
		{
			name:      "events today",
			date:      now,
			wantCount: 1,
		},
		{
			name:      "events in two days",
			date:      now.Add(48 * time.Hour),
			wantCount: 1,
		},
		{
			name:      "no events on date",
			date:      now.Add(96 * time.Hour),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsByDate(context.Background(), tt.date)
			if err != nil {
				t.Errorf("GetEventsByDate() error = %v", err)
				return
			}

			if len(events) != tt.wantCount {
				t.Errorf("GetEventsByDate() got %v events, want %v", len(events), tt.wantCount)
			}
		})
	}
}
