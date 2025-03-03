package astronomyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"astralis/internal/core/domain"
)

const (
	baseURL = "https://api.visibleplanets.dev/v3"
)

type astronomyAPIRepository struct {
	httpClient *http.Client
}

type planetVisibility struct {
	Name string `json:"name"`
	Constellation string  `json:"constellation"`
	Altitude      float64 `json:"altitude"`
	Azimuth       float64 `json:"azimuth"`
}

type visibilityResponse struct {
	Data []planetVisibility `json:"data"`
}

// NewAstronomyAPIRepository creates a new instance of the Astronomy API repository
func NewAstronomyAPIRepository() *astronomyAPIRepository {
	return &astronomyAPIRepository{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *astronomyAPIRepository) GetEvents(ctx context.Context, timeRange domain.TimeRange) ([]domain.Event, error) {
	// Get visible planets data
	url := fmt.Sprintf("%s?latitude=32&longitude=-98&date=%s",
		baseURL,
		timeRange.Start.Format("2006-01-02"),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching visible planets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("astronomy API returned status: %s", resp.Status)
	}

	var visResponse visibilityResponse
	if err := json.NewDecoder(resp.Body).Decode(&visResponse); err != nil {
		return nil, fmt.Errorf("decoding visibility response: %w", err)
	}

	var events []domain.Event
	for _, planet := range visResponse.Data {
		event := domain.Event{
			ID:          fmt.Sprintf("planet-%s-%s", planet.Name, timeRange.Start.Format("2006-01-02")),
			Title:       fmt.Sprintf("%s Visible in %s", planet.Name, planet.Constellation),
			Description: fmt.Sprintf("%s is visible at altitude %.2f° and azimuth %.2f°",
				planet.Name, planet.Altitude, planet.Azimuth),
			StartTime:   timeRange.Start,
			EndTime:     timeRange.Start.Add(24 * time.Hour),
			Type:        domain.Transit,
			Location:    planet.Constellation,
			Source:     "Visible Planets API",
			Visibility: fmt.Sprintf("Altitude: %.2f°, Azimuth: %.2f°", 
				planet.Altitude, planet.Azimuth),
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *astronomyAPIRepository) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	// For this API, we'll need to fetch all events and filter by ID
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now,
		End:   now.AddDate(0, 0, 7), // Look forward 1 week
	}

	events, err := r.GetEvents(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.ID == id {
			return &event, nil
		}
	}

	return nil, nil
}

func (r *astronomyAPIRepository) GetEventsByType(ctx context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
	events, err := r.GetEvents(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	var filtered []domain.Event
	for _, event := range events {
		if event.Type == eventType {
			filtered = append(filtered, event)
		}
	}

	return filtered, nil
}

func (r *astronomyAPIRepository) Name() string {
	return "Visible Planets API"
}
