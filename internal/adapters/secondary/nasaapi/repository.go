package nasaapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"astralis/internal/core/domain"
)

const (
	nasaBaseURL = "https://api.nasa.gov"
	cameBaseURL = "https://api.nasa.gov/DONKI"
)

type nasaAPIRepository struct {
	apiKey     string
	httpClient *http.Client
}

type cmeEvent struct {
	ActivityID     string    `json:"activityID"`
	StartTime      time.Time `json:"startTime"`
	Note           string    `json:"note"`
	CatalogID      string    `json:"catalog"`
	SourceLocation string    `json:"sourceLocation"`
}

// NewNASARepository creates a new NASA API repository
func NewNASARepository(apiKey string) *nasaAPIRepository {
	return &nasaAPIRepository{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *nasaAPIRepository) GetEvents(ctx context.Context, timeRange domain.TimeRange) ([]domain.Event, error) {
	// Get solar events (CMEs)
	url := fmt.Sprintf("%s/CME/?start_date=%s&end_date=%s&api_key=%s",
		cameBaseURL,
		timeRange.Start.Format("2006-01-02"),
		timeRange.End.Format("2006-01-02"),
		r.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching CME events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NASA API returned status: %s", resp.Status)
	}

	var rawEvents []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawEvents); err != nil {
		return nil, fmt.Errorf("decoding CME events: %w", err)
	}
	var cmeEvents []cmeEvent
	for _, raw := range rawEvents {
		var event cmeEvent
		event.ActivityID = raw["activityID"].(string)
		event.Note = raw["note"].(string)
		event.CatalogID = raw["catalog"].(string)
		event.SourceLocation = raw["sourceLocation"].(string)

		startTime, err := time.Parse("2006-01-02T15:04Z", raw["startTime"].(string))
		if err != nil {
			return nil, fmt.Errorf("parsing start time: %w", err)
		}
		event.StartTime = startTime
		cmeEvents = append(cmeEvents, event)
	}

	var events []domain.Event
	for _, cme := range cmeEvents {
		event := domain.Event{
			ID:          cme.ActivityID,
			Title:       fmt.Sprintf("Solar CME Event - %s", cme.SourceLocation),
			Description: cme.Note,
			StartTime:   cme.StartTime,
			EndTime:     cme.StartTime.Add(24 * time.Hour), // Approximate duration
			Type:        domain.Other,
			Location:    cme.SourceLocation,
			Source:      "NASA DONKI API",
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *nasaAPIRepository) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	// For NASA API, we'll need to fetch all events and filter by ID
	// This is not efficient but NASA's free API doesn't provide a direct lookup
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.AddDate(0, -1, 0), // Look back 1 month
		End:   now.AddDate(0, 1, 0),  // Look forward 1 month
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

func (r *nasaAPIRepository) GetEventsByType(ctx context.Context, eventType domain.EventType, timeRange domain.TimeRange) ([]domain.Event, error) {
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

func (r *nasaAPIRepository) Name() string {
	return "NASA API"
} 