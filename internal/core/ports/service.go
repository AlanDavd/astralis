package ports

import (
	"context"
	"time"

	"astralis/internal/core/domain"
)

// EventService defines the interface for the business logic layer
type EventService interface {
	// GetUpcomingEvents retrieves upcoming events within a specific time range
	GetUpcomingEvents(ctx context.Context, timeRange domain.TimeRange) ([]domain.Event, error)

	// GetEventByID retrieves a specific event by its ID
	GetEventByID(ctx context.Context, id string) (*domain.Event, error)

	// GetEventsByDate retrieves events for a specific date
	GetEventsByDate(ctx context.Context, date time.Time) ([]domain.Event, error)

	// GetEventsByDateRange retrieves events within a date range
	GetEventsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)

	// GetEventsByType retrieves events of a specific type
	GetEventsByType(ctx context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error)
}
