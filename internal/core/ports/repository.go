package ports

import (
	"context"

	"astralis/internal/core/domain"
)

// EventRepository defines the interface for accessing astronomical events
type EventRepository interface {
	// GetEvents retrieves events within a specific time range
	GetEvents(ctx context.Context, timeRange domain.TimeRange) ([]domain.Event, error)
	
	// GetEventByID retrieves a specific event by its ID
	GetEventByID(ctx context.Context, id string) (*domain.Event, error)
	
	// GetEventsByType retrieves events of a specific type within a time range
	GetEventsByType(ctx context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error)
	
	// Name returns the name of the repository implementation
	Name() string
}
