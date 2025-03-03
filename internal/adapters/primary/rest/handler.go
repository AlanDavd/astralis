package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"astralis/internal/core/domain"
	"astralis/internal/core/ports"
)

type Handler struct {
	service ports.EventService
}

func NewHandler(service ports.EventService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/events", h.GetEvents).Methods(http.MethodGet)
	router.HandleFunc("/events/{id}", h.GetEventByID).Methods(http.MethodGet)
	router.HandleFunc("/events/type/{type}", h.GetEventsByType).Methods(http.MethodGet)
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		start = time.Now()
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		end = start.AddDate(0, 1, 0) // Default to 1 month range
	}

	timeRange := domain.TimeRange{
		Start: start,
		End:   end,
	}

	events, err := h.service.GetUpcomingEvents(r.Context(), timeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	event, err := h.service.GetEventByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *Handler) GetEventsByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventType := domain.EventType(vars["type"])

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	} else {
		start = time.Now()
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	} else {
		end = start.AddDate(0, 1, 0) // Default to 1 month range
	}

	timeRange := domain.TimeRange{
		Start: start,
		End:   end,
	}

	events, err := h.service.GetEventsByType(r.Context(), eventType, timeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
} 