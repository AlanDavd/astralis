package domain

import (
	"time"
)

// EventType represents the type of astronomical event
type EventType string

const (
	MeteorShower EventType = "METEOR_SHOWER"
	Eclipse      EventType = "ECLIPSE"
	Conjunction  EventType = "CONJUNCTION"
	Transit      EventType = "TRANSIT"
	Other        EventType = "OTHER"
)

// Event represents an astronomical event
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Type        EventType `json:"type"`
	Visibility  string    `json:"visibility,omitempty"`
	Location    string    `json:"location,omitempty"`
	Source      string    `json:"source"`
}

// TimeRange represents a time period for filtering events
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// IsValid checks if the event has all required fields
func (e *Event) IsValid() bool {
	return e.Title != "" && e.Description != "" && !e.StartTime.IsZero()
}

// IsVisible checks if the event is visible at a given time
func (e *Event) IsVisible(t time.Time) bool {
	return (t.Equal(e.StartTime) || t.After(e.StartTime)) &&
		(e.EndTime.IsZero() || t.Equal(e.EndTime) || t.Before(e.EndTime))
}
