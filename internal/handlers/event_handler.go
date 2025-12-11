package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{service: svc}
}

// ListEvents handles GET /api/v1/events
func (h *EventHandler) ListEvents(c *gin.Context) {
	var req models.EventListRequest

	// Parse query parameters
	req.Namespace = c.Query("namespace")
	req.Type = c.Query("type")
	req.Since = c.Query("since")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	response, err := h.service.ListEvents(req)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to retrieve cluster events", err.Error())
		return
	}

	utils.ApiSuccess(c, response, "successfully retrieved cluster events")
}

// GetRecentEvents handles GET /api/v1/events/recent
func (h *EventHandler) GetRecentEvents(c *gin.Context) {
	limit := 10 // Default limit for recent events
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	events, err := h.service.GetRecentEvents(limit)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to retrieve recent events", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"events": events,
		"total":  len(events),
	}, "successfully retrieved recent events")
}

// GetEventsByObject handles GET /api/v1/events/object/:kind/:name
func (h *EventHandler) GetEventsByObject(c *gin.Context) {
	kind := c.Param("kind")
	name := c.Param("name")
	namespace := c.Query("namespace")

	if kind == "" || name == "" {
		utils.ApiError(c, http.StatusBadRequest, "invalid parameters", "kind and name are required")
		return
	}

	events, err := h.service.GetEventsByObject(namespace, kind, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to retrieve object events", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"events": events,
		"total":  len(events),
		"object": gin.H{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		},
	}, "successfully retrieved object events")
}
