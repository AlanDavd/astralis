package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"astralis/internal/core/domain"
)

// mockService implements ports.EventService for testing
type mockService struct {
	events map[string]domain.Event
}

func newMockService() *mockService {
	now := time.Now()
	events := map[string]domain.Event{
		"test-1": {
			ID:          "test-1",
			Title:       "Test Event 1",
			Description: "Test Description 1",
			StartTime:   now.Add(1 * time.Hour),
			EndTime:     now.Add(24 * time.Hour),
			Type:        domain.MeteorShower,
		},
		"test-2": {
			ID:          "test-2",
			Title:       "Test Event 2",
			Description: "Test Description 2",
			StartTime:   now.Add(2 * time.Hour),
			EndTime:     now.Add(48 * time.Hour),
			Type:        domain.Eclipse,
		},
	}

	return &mockService{events: events}
}

func (s *mockService) GetUpcomingEvents(_ context.Context, timeRange domain.TimeRange) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range s.events {
		if event.StartTime.After(timeRange.Start) && event.StartTime.Before(timeRange.End) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *mockService) GetEventByID(_ context.Context, id string) (*domain.Event, error) {
	if event, ok := s.events[id]; ok {
		return &event, nil
	}
	return nil, nil
}

func (s *mockService) GetEventsByDate(_ context.Context, date time.Time) ([]domain.Event, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	var result []domain.Event
	for _, event := range s.events {
		if event.StartTime.After(startOfDay) && event.StartTime.Before(endOfDay) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *mockService) GetEventsByDateRange(_ context.Context, start, end time.Time) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range s.events {
		if event.StartTime.After(start) && event.StartTime.Before(end) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *mockService) GetEventsByType(_ context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range s.events {
		if event.Type == eventType && event.StartTime.After(timeRange.Start) && event.StartTime.Before(timeRange.End) {
			result = append(result, event)
		}
	}
	return result, nil
}

func TestHandler_GetEvents(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetEvents() status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response []domain.Event
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("GetEvents() error decoding response = %v", err)
	}

	if len(response) != 2 {
		t.Errorf("GetEvents() response length = %v, want %v", len(response), 2)
	}
}

func TestHandler_GetEventByID(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	tests := []struct {
		name       string
		eventID    string
		wantStatus int
	}{
		{
			name:       "existing event",
			eventID:    "test-1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing event",
			eventID:    "non-existing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/"+tt.eventID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetEventByID() status code = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response domain.Event
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("GetEventByID() error decoding response = %v", err)
				}

				if response.ID != tt.eventID {
					t.Errorf("GetEventByID() response ID = %v, want %v", response.ID, tt.eventID)
				}
			}
		})
	}
}

func TestHandler_GetEventsByType(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	tests := []struct {
		name       string
		eventType  domain.EventType
		wantCount  int
		wantStatus int
	}{
		{
			name:       "meteor shower events",
			eventType:  domain.MeteorShower,
			wantCount:  1,
			wantStatus: http.StatusOK,
		},
		{
			name:       "eclipse events",
			eventType:  domain.Eclipse,
			wantCount:  1,
			wantStatus: http.StatusOK,
		},
		{
			name:       "no events of type",
			eventType:  domain.Transit,
			wantCount:  0,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/events/type/"+string(tt.eventType), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetEventsByType() status code = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response []domain.Event
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("GetEventsByType() error decoding response = %v", err)
				}

				if len(response) != tt.wantCount {
					t.Errorf("GetEventsByType() response length = %v, want %v", len(response), tt.wantCount)
				}

				for _, event := range response {
					if event.Type != tt.eventType {
						t.Errorf("GetEventsByType() event type = %v, want %v", event.Type, tt.eventType)
					}
				}
			}
		})
	}
} 