package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

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

func (s *mockService) GetEventsByType(_ context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
	var result []domain.Event
	for _, event := range s.events {
		if event.Type == eventType && event.StartTime.After(timeRange.Start) && event.StartTime.Before(timeRange.End) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *mockService) GetEventsByDate(_ context.Context, date time.Time) ([]domain.Event, error) {
	var result []domain.Event
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

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

func TestHandler_GetEvents(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := gin.Default()
	handler.RegisterRoutes(router)

	now := time.Now()
	startTime := now.Format(time.RFC3339)
	endTime := now.Add(96 * time.Hour).Format(time.RFC3339)

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "with time range",
			url:        fmt.Sprintf("/events?start=%s&end=%s", startTime, endTime),
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "without time range",
			url:        "/events",
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "invalid start date",
			url:        fmt.Sprintf("/events?start=invalid&end=%s", endTime),
			wantStatus: http.StatusBadRequest,
			wantCount:  0,
		},
		{
			name:       "invalid end date",
			url:        fmt.Sprintf("/events?start=%s&end=invalid", startTime),
			wantStatus: http.StatusBadRequest,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetEvents() status code = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response struct {
					Events []domain.Event `json:"events"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("GetEvents() error decoding response = %v", err)
				}

				if len(response.Events) != tt.wantCount {
					t.Errorf("GetEvents() response length = %v, want %v", len(response.Events), tt.wantCount)
				}
			}
		})
	}
}

func TestHandler_GetEventByID(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := gin.Default()
	handler.RegisterRoutes(router)

	tests := []struct {
		name       string
		eventID    string
		wantStatus int
		wantID     string
	}{
		{
			name:       "existing event",
			eventID:    "test-1",
			wantStatus: http.StatusOK,
			wantID:     "test-1",
		},
		{
			name:       "non-existing event",
			eventID:    "non-existing",
			wantStatus: http.StatusNotFound,
			wantID:     "",
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
				var response struct {
					Event domain.Event `json:"event"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("GetEventByID() error decoding response = %v", err)
				}

				if response.Event.ID != tt.wantID {
					t.Errorf("GetEventByID() response ID = %v, want %v", response.Event.ID, tt.wantID)
				}
			}
		})
	}
}

func TestHandler_GetEventsByType(t *testing.T) {
	mockSvc := newMockService()
	handler := NewHandler(mockSvc)

	router := gin.Default()
	handler.RegisterRoutes(router)

	now := time.Now()
	startTime := now.Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := now.Add(96 * time.Hour).Format(time.RFC3339)

	tests := []struct {
		name       string
		eventType  domain.EventType
		startTime  string
		endTime    string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "meteor shower events",
			eventType:  domain.MeteorShower,
			startTime:  startTime,
			endTime:    endTime,
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "eclipse events",
			eventType:  domain.Eclipse,
			startTime:  startTime,
			endTime:    endTime,
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "no events of type",
			eventType:  domain.Transit,
			startTime:  startTime,
			endTime:    endTime,
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "invalid start date",
			eventType:  domain.MeteorShower,
			startTime:  "invalid-date",
			endTime:    endTime,
			wantStatus: http.StatusBadRequest,
			wantCount:  0,
		},
		{
			name:       "invalid end date",
			eventType:  domain.MeteorShower,
			startTime:  startTime,
			endTime:    "invalid-date",
			wantStatus: http.StatusBadRequest,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/events/type/%s?start=%s&end=%s", 
				string(tt.eventType), tt.startTime, tt.endTime)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetEventsByType() status code = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response struct {
					Events []domain.Event `json:"events"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("GetEventsByType() error decoding response = %v", err)
				}

				if len(response.Events) != tt.wantCount {
					t.Errorf("GetEventsByType() response length = %v, want %v", len(response.Events), tt.wantCount)
				}

				for _, event := range response.Events {
					if event.Type != tt.eventType {
						t.Errorf("GetEventsByType() event type = %v, want %v", event.Type, tt.eventType)
					}
				}
			}
		})
	}
}
