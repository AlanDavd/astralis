package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/events", h.GetEvents)
	router.GET("/events/:id", h.GetEventByID)
	router.GET("/events/type/:type", h.GetEventsByType)
}

func (h *Handler) GetEvents(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	} else {
		start = time.Now()
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	} else {
		end = start.AddDate(0, 1, 0) // Default to 1 month range
	}

	timeRange := domain.TimeRange{
		Start: start,
		End:   end,
	}

	events, err := h.service.GetUpcomingEvents(c.Request.Context(), timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (h *Handler) GetEventByID(c *gin.Context) {
	id := c.Param("id")

	event, err := h.service.GetEventByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
	})
}

func (h *Handler) GetEventsByType(c *gin.Context) {
	eventType := domain.EventType(c.Param("type"))

	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	} else {
		start = time.Now()
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	} else {
		end = start.AddDate(0, 0, 1) // Default to 1 day range
	}

	timeRange := domain.TimeRange{
		Start: start,
		End:   end,
	}

	events, err := h.service.GetEventsByType(c.Request.Context(), eventType, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
} 