package handlers

import (
	"net/http"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// PrometheusHandler exposes Prometheus query APIs.
type PrometheusHandler struct {
	service *service.PrometheusService
}

func NewPrometheusHandler(svc *service.PrometheusService) *PrometheusHandler {
	return &PrometheusHandler{service: svc}
}

func (h *PrometheusHandler) GetStatus(c *gin.Context) {
	status, err := h.service.GetStatus(c.Request.Context())
	if err != nil {
		utils.ApiError(c, http.StatusBadGateway, "failed to get prometheus status", err.Error())
		return
	}
	utils.ApiSuccess(c, status, "prometheus status retrieved successfully")
}

func (h *PrometheusHandler) Query(c *gin.Context) {
	query := c.Query("query")
	var ts *time.Time
	if timeStr := c.Query("time"); timeStr != "" {
		parsed, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			utils.ApiError(c, http.StatusBadRequest, "invalid time, use RFC3339", err.Error())
			return
		}
		ts = &parsed
	}

	result, err := h.service.Query(c.Request.Context(), query, ts)
	if err != nil {
		utils.ApiError(c, http.StatusBadGateway, "prometheus query failed", err.Error())
		return
	}
	utils.ApiSuccess(c, result, "prometheus query succeeded")
}

func (h *PrometheusHandler) QueryRange(c *gin.Context) {
	query := c.Query("query")
	startStr := c.Query("start")
	endStr := c.Query("end")
	step := c.DefaultQuery("step", "60s")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "invalid start time, use RFC3339", err.Error())
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "invalid end time, use RFC3339", err.Error())
		return
	}

	result, err := h.service.QueryRange(c.Request.Context(), query, start, end, step)
	if err != nil {
		utils.ApiError(c, http.StatusBadGateway, "prometheus query_range failed", err.Error())
		return
	}
	utils.ApiSuccess(c, result, "prometheus query_range succeeded")
}
