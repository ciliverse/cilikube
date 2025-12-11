package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

// MonitoringService provides real-time monitoring and alerting
type MonitoringService struct {
	store        store.Store
	config       *configs.Config
	auditService *AuditService

	// Real-time metrics
	metrics      *RealTimeMetrics
	metricsMutex sync.RWMutex

	// Alert channels
	alertChannels []AlertChannel

	// Monitoring state
	isRunning bool
	stopChan  chan bool
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(store store.Store, config *configs.Config, auditService *AuditService) *MonitoringService {
	return &MonitoringService{
		store:         store,
		config:        config,
		auditService:  auditService,
		metrics:       NewRealTimeMetrics(),
		alertChannels: make([]AlertChannel, 0),
		stopChan:      make(chan bool),
	}
}

// RealTimeMetrics holds real-time security and system metrics
type RealTimeMetrics struct {
	// Authentication metrics
	LoginAttemptsPerMinute float64 `json:"login_attempts_per_minute"`
	FailedLoginsPerMinute  float64 `json:"failed_logins_per_minute"`
	ActiveSessions         int     `json:"active_sessions"`
	LockedAccounts         int     `json:"locked_accounts"`

	// Security metrics
	SecurityViolationsPerHour int `json:"security_violations_per_hour"`
	SuspiciousActivities      int `json:"suspicious_activities"`
	ActiveThreats             int `json:"active_threats"`

	// System metrics
	TotalUsers           int     `json:"total_users"`
	ActiveUsers          int     `json:"active_users"`
	APIRequestsPerMinute float64 `json:"api_requests_per_minute"`

	// Resource access metrics
	ResourceAccessPerMinute  float64 `json:"resource_access_per_minute"`
	PermissionDenialsPerHour int     `json:"permission_denials_per_hour"`

	// Timestamps
	LastUpdated    time.Time     `json:"last_updated"`
	UpdateInterval time.Duration `json:"update_interval"`
}

// NewRealTimeMetrics creates a new real-time metrics instance
func NewRealTimeMetrics() *RealTimeMetrics {
	return &RealTimeMetrics{
		LastUpdated:    time.Now(),
		UpdateInterval: 1 * time.Minute,
	}
}

// AlertLevel defines the severity of an alert
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelError    AlertLevel = "error"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert represents a security or system alert
type Alert struct {
	ID          string                 `json:"id"`
	Level       AlertLevel             `json:"level"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

// AlertChannel defines an interface for alert delivery
type AlertChannel interface {
	SendAlert(alert Alert) error
	GetName() string
}

// LogAlertChannel sends alerts to logs
type LogAlertChannel struct {
	name string
}

func NewLogAlertChannel() *LogAlertChannel {
	return &LogAlertChannel{name: "log"}
}

func (c *LogAlertChannel) SendAlert(alert Alert) error {
	fmt.Printf("[ALERT] %s - %s: %s\n", alert.Level, alert.Type, alert.Description)
	return nil
}

func (c *LogAlertChannel) GetName() string {
	return c.name
}

// Start begins the monitoring process
func (m *MonitoringService) Start() error {
	if m.isRunning {
		return fmt.Errorf("monitoring service is already running")
	}

	m.isRunning = true

	// Add default alert channel
	m.alertChannels = append(m.alertChannels, NewLogAlertChannel())

	// Start monitoring goroutines
	go m.metricsCollector()
	go m.threatDetector()
	go m.alertProcessor()

	return nil
}

// Stop stops the monitoring process
func (m *MonitoringService) Stop() error {
	if !m.isRunning {
		return fmt.Errorf("monitoring service is not running")
	}

	m.isRunning = false
	close(m.stopChan)

	return nil
}

// GetRealTimeMetrics returns current real-time metrics
func (m *MonitoringService) GetRealTimeMetrics() *RealTimeMetrics {
	m.metricsMutex.RLock()
	defer m.metricsMutex.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := *m.metrics
	return &metricsCopy
}

// metricsCollector collects and updates real-time metrics
func (m *MonitoringService) metricsCollector() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.updateMetrics()
		case <-m.stopChan:
			return
		}
	}
}

// updateMetrics updates the real-time metrics
func (m *MonitoringService) updateMetrics() {
	m.metricsMutex.Lock()
	defer m.metricsMutex.Unlock()

	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)

	// Get recent audit logs
	logs, _, err := m.store.ListAuditLogs(0, 1000)
	if err != nil {
		fmt.Printf("Error getting audit logs for metrics: %v\n", err)
		return
	}

	// Reset counters
	loginAttempts := 0
	failedLogins := 0
	securityViolations := 0
	apiRequests := 0
	resourceAccess := 0
	permissionDenials := 0

	// Count events in different time windows
	for _, log := range logs {
		// Count events in the last minute
		if log.CreatedAt.After(oneMinuteAgo) {
			switch log.Action {
			case "login", "login_failed":
				loginAttempts++
				if log.Action == "login_failed" {
					failedLogins++
				}
			case "resource_access":
				resourceAccess++
			}
			apiRequests++
		}

		// Count events in the last hour
		if log.CreatedAt.After(oneHourAgo) {
			switch log.Action {
			case "permission_denied":
				permissionDenials++
			case "security_violation", "suspicious_activity":
				securityViolations++
			}
		}
	}

	// Update metrics
	m.metrics.LoginAttemptsPerMinute = float64(loginAttempts)
	m.metrics.FailedLoginsPerMinute = float64(failedLogins)
	m.metrics.SecurityViolationsPerHour = securityViolations
	m.metrics.APIRequestsPerMinute = float64(apiRequests)
	m.metrics.ResourceAccessPerMinute = float64(resourceAccess)
	m.metrics.PermissionDenialsPerHour = permissionDenials

	// Get user statistics
	users, total, err := m.store.ListUsers(0, 10000)
	if err == nil {
		m.metrics.TotalUsers = int(total)
		activeUsers := 0
		for _, user := range users {
			if user.IsActive {
				activeUsers++
			}
		}
		m.metrics.ActiveUsers = activeUsers
	}

	// Update timestamp
	m.metrics.LastUpdated = now

	// Check for alert conditions
	m.checkAlertConditions()
}

// checkAlertConditions checks if any alert conditions are met
func (m *MonitoringService) checkAlertConditions() {
	// Check for high failed login rate
	if m.metrics.FailedLoginsPerMinute > 10 {
		m.createAlert(AlertLevelWarning, "high_failed_logins",
			"High Failed Login Rate",
			fmt.Sprintf("%.0f failed logins per minute detected", m.metrics.FailedLoginsPerMinute),
			map[string]interface{}{
				"rate":      m.metrics.FailedLoginsPerMinute,
				"threshold": 10,
			})
	}

	// Check for high permission denials
	if m.metrics.PermissionDenialsPerHour > 50 {
		m.createAlert(AlertLevelWarning, "high_permission_denials",
			"High Permission Denials",
			fmt.Sprintf("%d permission denials per hour detected", m.metrics.PermissionDenialsPerHour),
			map[string]interface{}{
				"count":     m.metrics.PermissionDenialsPerHour,
				"threshold": 50,
			})
	}

	// Check for security violations
	if m.metrics.SecurityViolationsPerHour > 5 {
		m.createAlert(AlertLevelError, "security_violations",
			"Security Violations Detected",
			fmt.Sprintf("%d security violations per hour detected", m.metrics.SecurityViolationsPerHour),
			map[string]interface{}{
				"count":     m.metrics.SecurityViolationsPerHour,
				"threshold": 5,
			})
	}
}

// threatDetector runs threat detection periodically
func (m *MonitoringService) threatDetector() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.detectThreats()
		case <-m.stopChan:
			return
		}
	}
}

// detectThreats detects security threats and creates alerts
func (m *MonitoringService) detectThreats() {
	threats, err := m.auditService.DetectAnomalousActivity()
	if err != nil {
		fmt.Printf("Error detecting threats: %v\n", err)
		return
	}

	m.metricsMutex.Lock()
	m.metrics.ActiveThreats = len(threats)
	m.metricsMutex.Unlock()

	// Create alerts for detected threats
	for _, threat := range threats {
		alertLevel := AlertLevelWarning
		if threat.Severity == SeverityError || threat.Severity == SeverityCritical {
			alertLevel = AlertLevelError
		}

		m.createAlert(alertLevel, threat.Type,
			fmt.Sprintf("Security Threat: %s", threat.Type),
			threat.Description,
			map[string]interface{}{
				"threat_type": threat.Type,
				"severity":    threat.Severity,
				"count":       threat.Count,
				"ip_address":  threat.IPAddress,
				"user_id":     threat.UserID,
				"first_seen":  threat.FirstSeen,
				"last_seen":   threat.LastSeen,
				"details":     threat.Details,
			})
	}
}

// alertProcessor processes and sends alerts
func (m *MonitoringService) alertProcessor() {
	// This would typically read from an alert queue
	// For now, it's a placeholder for alert processing logic
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Process any pending alerts
			// In a real implementation, this would handle alert deduplication,
			// rate limiting, and delivery to various channels
		case <-m.stopChan:
			return
		}
	}
}

// createAlert creates and sends an alert
func (m *MonitoringService) createAlert(level AlertLevel, alertType, title, description string, data map[string]interface{}) {
	alert := Alert{
		ID:          fmt.Sprintf("%s_%d", alertType, time.Now().Unix()),
		Level:       level,
		Type:        alertType,
		Title:       title,
		Description: description,
		Source:      "monitoring_service",
		Timestamp:   time.Now(),
		Data:        data,
		Resolved:    false,
	}

	// Send alert through all channels
	for _, channel := range m.alertChannels {
		if err := channel.SendAlert(alert); err != nil {
			fmt.Printf("Error sending alert through channel %s: %v\n", channel.GetName(), err)
		}
	}

	// Log the alert as a security event
	m.auditService.LogSecurityEvent(SecurityEvent{
		Type:      "alert_generated",
		Severity:  string(level),
		Resource:  "monitoring",
		Action:    "alert",
		Result:    "generated",
		Details:   data,
		Timestamp: time.Now(),
	})
}

// GetSystemHealth returns overall system health status
func (m *MonitoringService) GetSystemHealth() *SystemHealth {
	metrics := m.GetRealTimeMetrics()

	health := &SystemHealth{
		Status:    "healthy",
		Timestamp: time.Now(),
		Metrics:   metrics,
		Issues:    make([]HealthIssue, 0),
	}

	// Check for health issues
	if metrics.FailedLoginsPerMinute > 20 {
		health.Status = "warning"
		health.Issues = append(health.Issues, HealthIssue{
			Type:        "security",
			Severity:    "warning",
			Description: "High rate of failed login attempts",
			Value:       metrics.FailedLoginsPerMinute,
			Threshold:   20,
		})
	}

	if metrics.ActiveThreats > 0 {
		health.Status = "critical"
		health.Issues = append(health.Issues, HealthIssue{
			Type:        "security",
			Severity:    "critical",
			Description: "Active security threats detected",
			Value:       float64(metrics.ActiveThreats),
			Threshold:   0,
		})
	}

	if metrics.SecurityViolationsPerHour > 10 {
		if health.Status == "healthy" {
			health.Status = "warning"
		}
		health.Issues = append(health.Issues, HealthIssue{
			Type:        "security",
			Severity:    "warning",
			Description: "High rate of security violations",
			Value:       float64(metrics.SecurityViolationsPerHour),
			Threshold:   10,
		})
	}

	return health
}

// SystemHealth represents the overall system health
type SystemHealth struct {
	Status    string           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Metrics   *RealTimeMetrics `json:"metrics"`
	Issues    []HealthIssue    `json:"issues"`
}

// HealthIssue represents a system health issue
type HealthIssue struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	Threshold   float64 `json:"threshold"`
}

// AddAlertChannel adds a new alert channel
func (m *MonitoringService) AddAlertChannel(channel AlertChannel) {
	m.alertChannels = append(m.alertChannels, channel)
}
