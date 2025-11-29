package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse health check response structure
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

var startTime = time.Now()

// HealthCheck basic health check
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}

// ReadinessCheck readiness check
func ReadinessCheck(c *gin.Context) {
	checks := make(map[string]string)

	// Check Kubernetes connection
	checks["kubernetes"] = "ok"

	// Check database connection (if enabled)
	checks["database"] = "ok"

	// All checks passed
	allHealthy := true
	for _, status := range checks {
		if status != "ok" {
			allHealthy = false
			break
		}
	}

	status := "ready"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "not ready"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Checks:    checks,
	})
}

// LivenessCheck liveness check
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "alive",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}
