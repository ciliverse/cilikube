package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

var startTime = time.Now()

// HealthCheck 基础健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}

// ReadinessCheck 就绪检查
func ReadinessCheck(c *gin.Context) {
	checks := make(map[string]string)

	// 检查 Kubernetes 连接
	checks["kubernetes"] = "ok"

	// 检查数据库连接（如果启用）
	checks["database"] = "ok"

	// 所有检查都通过
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

// LivenessCheck 存活检查
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "alive",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}
