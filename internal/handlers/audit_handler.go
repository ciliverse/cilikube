package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// GetAuditLogs gets audit logs with pagination and filtering
// @Summary Get audit logs
// @Description Get audit logs with optional filtering by user, action, and time range
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param user_id query int false "Filter by user ID"
// @Param action query string false "Filter by action"
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/logs [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userIDStr := c.Query("user_id")
	action := c.Query("action")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Parse time filters (for future use in filtering)
	if startTimeStr != "" {
		if _, parseErr := time.Parse(time.RFC3339, startTimeStr); parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid start_time format. Use RFC3339 format.",
			})
			return
		}
	}
	if endTimeStr != "" {
		if _, parseErr := time.Parse(time.RFC3339, endTimeStr); parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid end_time format. Use RFC3339 format.",
			})
			return
		}
	}

	// Get audit logs based on filters
	var logs interface{}
	var total int64
	var err error

	if userIDStr != "" {
		userID, parseErr := strconv.ParseUint(userIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid user_id format",
			})
			return
		}
		logs, total, err = h.auditService.GetAuditLogsByUserID(uint(userID), offset, pageSize)
	} else if action != "" {
		logs, total, err = h.auditService.GetAuditLogsByAction(action, offset, pageSize)
	} else {
		logs, total, err = h.auditService.GetAllAuditLogs(offset, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get audit logs: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Retrieved successfully",
		"data": gin.H{
			"logs":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAuditReport generates an audit report for a specific time period
// @Summary Get audit report
// @Description Generate comprehensive audit report for specified time period
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_time query string true "Start time (RFC3339 format)"
// @Param end_time query string true "End time (RFC3339 format)"
// @Param user_id query int false "Filter by user ID"
// @Success 200 {object} service.AuditReport
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/report [get]
func (h *AuditHandler) GetAuditReport(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	userIDStr := c.Query("user_id")

	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "start_time and end_time are required",
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid start_time format. Use RFC3339 format.",
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid end_time format. Use RFC3339 format.",
		})
		return
	}

	var userID *uint
	if userIDStr != "" {
		uid, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid user_id format",
			})
			return
		}
		uidUint := uint(uid)
		userID = &uidUint
	}

	report, err := h.auditService.GetAuditReport(startTime, endTime, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to generate audit report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Report generated successfully",
		"data":    report,
	})
}

// GetSecurityMetrics gets security metrics for monitoring
// @Summary Get security metrics
// @Description Get security metrics for the specified time period
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param period query string false "Time period (e.g., '24h', '7d', '30d')" default("24h")
// @Success 200 {object} service.SecurityMetrics
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/metrics [get]
func (h *AuditHandler) GetSecurityMetrics(c *gin.Context) {
	periodStr := c.DefaultQuery("period", "24h")

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid period format. Use duration format like '24h', '7d', etc.",
		})
		return
	}

	metrics, err := h.auditService.GetSecurityMetrics(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get security metrics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Metrics retrieved successfully",
		"data":    metrics,
	})
}

// DetectThreats detects security threats and anomalous activities
// @Summary Detect security threats
// @Description Analyze audit logs to detect potential security threats
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/threats [get]
func (h *AuditHandler) DetectThreats(c *gin.Context) {
	threats, err := h.auditService.DetectAnomalousActivity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to detect threats: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Threat detection completed",
		"data": gin.H{
			"threats": threats,
			"count":   len(threats),
		},
	})
}

// GetUserActivity gets activity summary for a specific user
// @Summary Get user activity
// @Description Get detailed activity summary for a specific user
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Param period query string false "Time period (e.g., '24h', '7d', '30d')" default("7d")
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/users/{user_id}/activity [get]
func (h *AuditHandler) GetUserActivity(c *gin.Context) {
	userIDStr := c.Param("user_id")
	periodStr := c.DefaultQuery("period", "7d")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid user_id format",
		})
		return
	}

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid period format. Use duration format like '24h', '7d', etc.",
		})
		return
	}

	startTime := time.Now().Add(-period)
	endTime := time.Now()
	uid := uint(userID)

	report, err := h.auditService.GetAuditReport(startTime, endTime, &uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get user activity: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User activity retrieved successfully",
		"data":    report,
	})
}

// GetSystemActivity gets overall system activity summary
// @Summary Get system activity
// @Description Get overall system activity and statistics
// @Tags Audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param period query string false "Time period (e.g., '24h', '7d', '30d')" default("24h")
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/audit/system/activity [get]
func (h *AuditHandler) GetSystemActivity(c *gin.Context) {
	periodStr := c.DefaultQuery("period", "24h")

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid period format. Use duration format like '24h', '7d', etc.",
		})
		return
	}

	startTime := time.Now().Add(-period)
	endTime := time.Now()

	// Get system-wide report
	report, err := h.auditService.GetAuditReport(startTime, endTime, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get system activity: " + err.Error(),
		})
		return
	}

	// Get security metrics
	metrics, err := h.auditService.GetSecurityMetrics(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get security metrics: " + err.Error(),
		})
		return
	}

	// Detect current threats
	threats, err := h.auditService.DetectAnomalousActivity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to detect threats: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "System activity retrieved successfully",
		"data": gin.H{
			"report":         report,
			"metrics":        metrics,
			"active_threats": threats,
			"summary": gin.H{
				"period":             periodStr,
				"total_events":       report.TotalEvents,
				"active_threats":     len(threats),
				"login_success_rate": report.LoginSuccessRate,
			},
		},
	})
}

// Helper methods for the audit handler

func (h *AuditHandler) GetAuditLogsByUserID(userID uint, offset, limit int) (interface{}, int64, error) {
	return h.auditService.GetAuditLogsByUserID(userID, offset, limit)
}

func (h *AuditHandler) GetAuditLogsByAction(action string, offset, limit int) (interface{}, int64, error) {
	return h.auditService.GetAuditLogsByAction(action, offset, limit)
}

func (h *AuditHandler) GetAllAuditLogs(offset, limit int) (interface{}, int64, error) {
	return h.auditService.GetAllAuditLogs(offset, limit)
}
