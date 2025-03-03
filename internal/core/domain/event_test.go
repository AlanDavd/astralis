package domain

import (
	"testing"
	"time"
)

func TestEvent_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		want  bool
	}{
		{
			name: "valid event",
			event: Event{
				Title:       "Test Event",
				Description: "Test Description",
				StartTime:   time.Now(),
			},
			want: true,
		},
		{
			name: "missing title",
			event: Event{
				Description: "Test Description",
				StartTime:   time.Now(),
			},
			want: false,
		},
		{
			name: "missing description",
			event: Event{
				Title:     "Test Event",
				StartTime: time.Now(),
			},
			want: false,
		},
		{
			name: "missing start time",
			event: Event{
				Title:       "Test Event",
				Description: "Test Description",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsValid(); got != tt.want {
				t.Errorf("Event.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_IsVisible(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		event     Event
		checkTime time.Time
		want      bool
	}{
		{
			name: "event is currently visible",
			event: Event{
				StartTime: now.Add(-1 * time.Hour),
				EndTime:   now.Add(1 * time.Hour),
			},
			checkTime: now,
			want:      true,
		},
		{
			name: "event hasn't started",
			event: Event{
				StartTime: now.Add(1 * time.Hour),
				EndTime:   now.Add(2 * time.Hour),
			},
			checkTime: now,
			want:      false,
		},
		{
			name: "event has ended",
			event: Event{
				StartTime: now.Add(-2 * time.Hour),
				EndTime:   now.Add(-1 * time.Hour),
			},
			checkTime: now,
			want:      false,
		},
		{
			name: "event with no end time",
			event: Event{
				StartTime: now.Add(-1 * time.Hour),
			},
			checkTime: now,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsVisible(tt.checkTime); got != tt.want {
				t.Errorf("Event.IsVisible() = %v, want %v", got, tt.want)
			}
		})
	}
} 