package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

type MonitoringHandler struct {
	monitoringService *service.MonitoringService
}

func NewMonitoringHandler(monitoringService *service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringService: monitoringService,
	}
}

// GetRealTimeMetrics gets real-time security and system metrics
// @Summary Get real-time metrics
// @Description Get current real-time security and system metrics
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.RealTimeMetrics
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/metrics [get]
func (h *MonitoringHandler) GetRealTimeMetrics(c *gin.Context) {
	metrics := h.monitoringService.GetRealTimeMetrics()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Metrics retrieved successfully",
		"data":    metrics,
	})
}

// GetSystemHealth gets overall system health status
// @Summary Get system health
// @Description Get overall system health status and issues
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.SystemHealth
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/health [get]
func (h *MonitoringHandler) GetSystemHealth(c *gin.Context) {
	health := h.monitoringService.GetSystemHealth()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "System health retrieved successfully",
		"data":    health,
	})
}

// GetDashboardData gets comprehensive dashboard data
// @Summary Get dashboard data
// @Description Get comprehensive monitoring dashboard data
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/dashboard [get]
func (h *MonitoringHandler) GetDashboardData(c *gin.Context) {
	// Get real-time metrics
	metrics := h.monitoringService.GetRealTimeMetrics()

	// Get system health
	health := h.monitoringService.GetSystemHealth()

	// Prepare dashboard data
	dashboardData := gin.H{
		"metrics": metrics,
		"health":  health,
		"summary": gin.H{
			"status":              health.Status,
			"total_users":         metrics.TotalUsers,
			"active_users":        metrics.ActiveUsers,
			"active_sessions":     metrics.ActiveSessions,
			"active_threats":      metrics.ActiveThreats,
			"failed_logins_rate":  metrics.FailedLoginsPerMinute,
			"security_violations": metrics.SecurityViolationsPerHour,
			"last_updated":        metrics.LastUpdated,
		},
		"alerts": gin.H{
			"critical": countIssuesBySeverity(health.Issues, "critical"),
			"warning":  countIssuesBySeverity(health.Issues, "warning"),
			"info":     countIssuesBySeverity(health.Issues, "info"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Dashboard data retrieved successfully",
		"data":    dashboardData,
	})
}

// GetMetricsHistory gets historical metrics data
// @Summary Get metrics history
// @Description Get historical metrics data for charts and trends
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param period query string false "Time period (e.g., '1h', '24h', '7d')" default("24h")
// @Param interval query string false "Data interval (e.g., '1m', '5m', '1h')" default("5m")
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/metrics/history [get]
func (h *MonitoringHandler) GetMetricsHistory(c *gin.Context) {
	periodStr := c.DefaultQuery("period", "24h")
	intervalStr := c.DefaultQuery("interval", "5m")

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid period format. Use duration format like '1h', '24h', etc.",
		})
		return
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid interval format. Use duration format like '1m', '5m', etc.",
		})
		return
	}

	// Generate historical data (in a real implementation, this would come from a time-series database)
	historyData := h.generateMetricsHistory(period, interval)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Metrics history retrieved successfully",
		"data":    historyData,
	})
}

// GetSecurityOverview gets security-focused overview
// @Summary Get security overview
// @Description Get security-focused monitoring overview
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/security [get]
func (h *MonitoringHandler) GetSecurityOverview(c *gin.Context) {
	metrics := h.monitoringService.GetRealTimeMetrics()
	health := h.monitoringService.GetSystemHealth()

	// Calculate security score (0-100)
	securityScore := h.calculateSecurityScore(metrics, health)

	securityOverview := gin.H{
		"security_score": securityScore,
		"threat_level":   h.getThreatLevel(metrics, health),
		"metrics": gin.H{
			"failed_logins_per_minute":     metrics.FailedLoginsPerMinute,
			"security_violations_per_hour": metrics.SecurityViolationsPerHour,
			"active_threats":               metrics.ActiveThreats,
			"locked_accounts":              metrics.LockedAccounts,
			"permission_denials_per_hour":  metrics.PermissionDenialsPerHour,
		},
		"recent_threats":  h.getRecentThreats(),
		"security_trends": h.getSecurityTrends(),
		"recommendations": h.getSecurityRecommendations(metrics, health),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Security overview retrieved successfully",
		"data":    securityOverview,
	})
}

// GetAlerts gets current system alerts
// @Summary Get system alerts
// @Description Get current system alerts and notifications
// @Tags Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param severity query string false "Filter by severity (info, warning, error, critical)"
// @Param limit query int false "Limit number of results" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/monitoring/alerts [get]
func (h *MonitoringHandler) GetAlerts(c *gin.Context) {
	severity := c.Query("severity")
	limitStr := c.DefaultQuery("limit", "50")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	// In a real implementation, this would fetch alerts from a database or queue
	alerts := h.generateSampleAlerts(severity, limit)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Alerts retrieved successfully",
		"data": gin.H{
			"alerts": alerts,
			"count":  len(alerts),
			"filter": gin.H{
				"severity": severity,
				"limit":    limit,
			},
		},
	})
}

// Helper methods

func countIssuesBySeverity(issues []service.HealthIssue, severity string) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

func (h *MonitoringHandler) generateMetricsHistory(period, interval time.Duration) map[string]interface{} {
	// This is a simplified implementation for demonstration
	// In a real system, this would query a time-series database

	now := time.Now()
	points := int(period / interval)
	if points > 1000 {
		points = 1000 // Limit data points
	}

	timestamps := make([]time.Time, points)
	loginAttempts := make([]float64, points)
	failedLogins := make([]float64, points)
	apiRequests := make([]float64, points)

	for i := 0; i < points; i++ {
		timestamp := now.Add(-period + time.Duration(i)*interval)
		timestamps[i] = timestamp

		// Generate sample data (in reality, this would come from stored metrics)
		loginAttempts[i] = float64(5 + i%10)
		failedLogins[i] = float64(1 + i%3)
		apiRequests[i] = float64(50 + i%20)
	}

	return map[string]interface{}{
		"period":     period.String(),
		"interval":   interval.String(),
		"timestamps": timestamps,
		"series": map[string]interface{}{
			"login_attempts": loginAttempts,
			"failed_logins":  failedLogins,
			"api_requests":   apiRequests,
		},
	}
}

func (h *MonitoringHandler) calculateSecurityScore(metrics *service.RealTimeMetrics, health *service.SystemHealth) int {
	score := 100

	// Deduct points for security issues
	if metrics.FailedLoginsPerMinute > 10 {
		score -= 20
	} else if metrics.FailedLoginsPerMinute > 5 {
		score -= 10
	}

	if metrics.ActiveThreats > 0 {
		score -= 30
	}

	if metrics.SecurityViolationsPerHour > 10 {
		score -= 25
	} else if metrics.SecurityViolationsPerHour > 5 {
		score -= 15
	}

	if metrics.PermissionDenialsPerHour > 50 {
		score -= 15
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

func (h *MonitoringHandler) getThreatLevel(metrics *service.RealTimeMetrics, health *service.SystemHealth) string {
	if metrics.ActiveThreats > 0 || health.Status == "critical" {
		return "high"
	}

	if metrics.FailedLoginsPerMinute > 15 || metrics.SecurityViolationsPerHour > 10 {
		return "medium"
	}

	return "low"
}

func (h *MonitoringHandler) getRecentThreats() []map[string]interface{} {
	// In a real implementation, this would fetch recent threats from the audit service
	return []map[string]interface{}{
		{
			"type":        "brute_force_login",
			"severity":    "warning",
			"description": "Multiple failed login attempts from IP 192.168.1.100",
			"timestamp":   time.Now().Add(-10 * time.Minute),
		},
	}
}

func (h *MonitoringHandler) getSecurityTrends() map[string]interface{} {
	// In a real implementation, this would calculate trends from historical data
	return map[string]interface{}{
		"failed_logins_trend":       "stable",
		"security_violations_trend": "decreasing",
		"threat_detection_trend":    "stable",
	}
}

func (h *MonitoringHandler) getSecurityRecommendations(metrics *service.RealTimeMetrics, health *service.SystemHealth) []string {
	recommendations := make([]string, 0)

	if metrics.FailedLoginsPerMinute > 10 {
		recommendations = append(recommendations, "Consider implementing additional rate limiting for login attempts")
	}

	if metrics.ActiveThreats > 0 {
		recommendations = append(recommendations, "Review and address active security threats immediately")
	}

	if metrics.PermissionDenialsPerHour > 50 {
		recommendations = append(recommendations, "Review user permissions and role assignments")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Security posture is good. Continue monitoring.")
	}

	return recommendations
}

func (h *MonitoringHandler) generateSampleAlerts(severity string, limit int) []map[string]interface{} {
	// In a real implementation, this would fetch alerts from a database
	alerts := make([]map[string]interface{}, 0)

	sampleAlerts := []map[string]interface{}{
		{
			"id":          "alert_001",
			"level":       "warning",
			"type":        "high_failed_logins",
			"title":       "High Failed Login Rate",
			"description": "15 failed logins per minute detected",
			"timestamp":   time.Now().Add(-5 * time.Minute),
			"resolved":    false,
		},
		{
			"id":          "alert_002",
			"level":       "error",
			"type":        "security_violation",
			"title":       "Security Violation",
			"description": "Unauthorized access attempt detected",
			"timestamp":   time.Now().Add(-15 * time.Minute),
			"resolved":    false,
		},
	}

	for _, alert := range sampleAlerts {
		if severity == "" || alert["level"] == severity {
			alerts = append(alerts, alert)
			if len(alerts) >= limit {
				break
			}
		}
	}

	return alerts
}
