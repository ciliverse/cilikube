package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

// AuditService provides audit and monitoring functionality
type AuditService struct {
	store  store.Store
	config *configs.Config
}

// NewAuditService creates a new AuditService instance
func NewAuditService(store store.Store, config *configs.Config) *AuditService {
	return &AuditService{
		store:  store,
		config: config,
	}
}

// SecurityEvent represents a security-related event
type SecurityEvent struct {
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"`
	UserID    *uint                  `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Resource  string                 `json:"resource"`
	Action    string                 `json:"action"`
	Result    string                 `json:"result"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
}

// AuditEventType defines types of audit events
type AuditEventType string

const (
	// Authentication events
	EventTypeLogin          AuditEventType = "login"
	EventTypeLoginFailed    AuditEventType = "login_failed"
	EventTypeLogout         AuditEventType = "logout"
	EventTypePasswordChange AuditEventType = "password_change"
	EventTypeAccountLocked  AuditEventType = "account_locked"

	// Authorization events
	EventTypePermissionDenied AuditEventType = "permission_denied"
	EventTypeRoleAssigned     AuditEventType = "role_assigned"
	EventTypeRoleRemoved      AuditEventType = "role_removed"

	// Resource access events
	EventTypeResourceAccess AuditEventType = "resource_access"
	EventTypeResourceCreate AuditEventType = "resource_create"
	EventTypeResourceUpdate AuditEventType = "resource_update"
	EventTypeResourceDelete AuditEventType = "resource_delete"

	// System events
	EventTypeSystemConfig AuditEventType = "system_config"
	EventTypeUserManage   AuditEventType = "user_manage"

	// Security events
	EventTypeSuspiciousActivity AuditEventType = "suspicious_activity"
	EventTypeSecurityViolation  AuditEventType = "security_violation"
	EventTypeRateLimitExceeded  AuditEventType = "rate_limit_exceeded"
)

// EventSeverity defines severity levels for events
type EventSeverity string

const (
	SeverityInfo     EventSeverity = "info"
	SeverityWarning  EventSeverity = "warning"
	SeverityError    EventSeverity = "error"
	SeverityCritical EventSeverity = "critical"
)

// LogSecurityEvent logs a security event
func (s *AuditService) LogSecurityEvent(event SecurityEvent) error {
	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Serialize details to JSON
	detailsJSON := ""
	if event.Details != nil {
		if jsonBytes, err := json.Marshal(event.Details); err == nil {
			detailsJSON = string(jsonBytes)
		}
	}

	// Create audit log entry
	auditLog := &store.AuditLog{
		UserID:     event.UserID,
		Action:     string(event.Type),
		Resource:   event.Resource,
		ResourceID: event.Action,
		IPAddress:  event.IPAddress,
		UserAgent:  event.UserAgent,
		Details:    detailsJSON,
		CreatedAt:  event.Timestamp,
	}

	return s.store.CreateAuditLog(auditLog)
}

// LogAuthenticationEvent logs authentication-related events
func (s *AuditService) LogAuthenticationEvent(eventType AuditEventType, userID *uint, username, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	severity := SeverityInfo
	if !success {
		severity = SeverityWarning
	}

	result := "success"
	if !success {
		result = "failure"
	}

	event := SecurityEvent{
		Type:      string(eventType),
		Severity:  string(severity),
		UserID:    userID,
		Username:  username,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Resource:  "authentication",
		Action:    string(eventType),
		Result:    result,
		Details:   details,
		Timestamp: time.Now(),
	}

	return s.LogSecurityEvent(event)
}

// LogResourceAccessEvent logs resource access events
func (s *AuditService) LogResourceAccessEvent(userID uint, username, resource, action, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	severity := SeverityInfo
	if !success {
		severity = SeverityWarning
	}

	result := "success"
	if !success {
		result = "failure"
	}

	event := SecurityEvent{
		Type:      string(EventTypeResourceAccess),
		Severity:  string(severity),
		UserID:    &userID,
		Username:  username,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Resource:  resource,
		Action:    action,
		Result:    result,
		Details:   details,
		Timestamp: time.Now(),
	}

	return s.LogSecurityEvent(event)
}

// DetectAnomalousActivity analyzes audit logs to detect suspicious patterns
func (s *AuditService) DetectAnomalousActivity() ([]SecurityThreat, error) {
	var threats []SecurityThreat

	// Detect multiple failed logins from same IP
	if ipThreats, err := s.detectFailedLoginsByIP(); err == nil {
		threats = append(threats, ipThreats...)
	}

	// Detect unusual access patterns
	if accessThreats, err := s.detectUnusualAccessPatterns(); err == nil {
		threats = append(threats, accessThreats...)
	}

	// Detect privilege escalation attempts
	if privThreats, err := s.detectPrivilegeEscalation(); err == nil {
		threats = append(threats, privThreats...)
	}

	// Detect brute force attacks
	if bruteForceThreats, err := s.detectBruteForceAttacks(); err == nil {
		threats = append(threats, bruteForceThreats...)
	}

	return threats, nil
}

// SecurityThreat represents a detected security threat
type SecurityThreat struct {
	Type        string                 `json:"type"`
	Severity    EventSeverity          `json:"severity"`
	Description string                 `json:"description"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserID      *uint                  `json:"user_id,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Count       int                    `json:"count"`
	FirstSeen   time.Time              `json:"first_seen"`
	LastSeen    time.Time              `json:"last_seen"`
	Details     map[string]interface{} `json:"details"`
}

// detectFailedLoginsByIP detects multiple failed logins from the same IP
func (s *AuditService) detectFailedLoginsByIP() ([]SecurityThreat, error) {
	var threats []SecurityThreat

	// Get failed login attempts from the last hour
	since := time.Now().Add(-1 * time.Hour)
	logs, _, err := s.store.GetAuditLogsByAction("login_failed", 0, 1000)
	if err != nil {
		return threats, err
	}

	// Group by IP address
	ipFailures := make(map[string][]store.AuditLog)
	for _, log := range logs {
		if log.CreatedAt.After(since) && log.IPAddress != "" {
			ipFailures[log.IPAddress] = append(ipFailures[log.IPAddress], *log)
		}
	}

	// Check for suspicious patterns
	for ip, failures := range ipFailures {
		if len(failures) >= 10 { // 10 or more failures from same IP
			threat := SecurityThreat{
				Type:        "brute_force_login",
				Severity:    SeverityError,
				Description: fmt.Sprintf("Multiple failed login attempts (%d) from IP %s", len(failures), ip),
				IPAddress:   ip,
				Count:       len(failures),
				FirstSeen:   failures[len(failures)-1].CreatedAt, // Oldest
				LastSeen:    failures[0].CreatedAt,               // Newest
				Details: map[string]interface{}{
					"failed_attempts": len(failures),
					"time_window":     "1 hour",
				},
			}
			threats = append(threats, threat)
		}
	}

	return threats, nil
}

// detectUnusualAccessPatterns detects unusual access patterns
func (s *AuditService) detectUnusualAccessPatterns() ([]SecurityThreat, error) {
	var threats []SecurityThreat

	// Get recent audit logs
	since := time.Now().Add(-24 * time.Hour)
	logs, _, err := s.store.ListAuditLogs(0, 1000)
	if err != nil {
		return threats, err
	}

	// Group by user and analyze patterns
	userActivity := make(map[uint][]store.AuditLog)
	for _, log := range logs {
		if log.CreatedAt.After(since) && log.UserID != nil {
			userActivity[*log.UserID] = append(userActivity[*log.UserID], *log)
		}
	}

	for userID, activities := range userActivity {
		// Check for unusual time patterns (access outside normal hours)
		offHoursCount := 0
		uniqueIPs := make(map[string]bool)

		for _, activity := range activities {
			hour := activity.CreatedAt.Hour()
			// Consider 22:00-06:00 as off-hours
			if hour >= 22 || hour <= 6 {
				offHoursCount++
			}
			if activity.IPAddress != "" {
				uniqueIPs[activity.IPAddress] = true
			}
		}

		// Alert if significant off-hours activity
		if offHoursCount > 5 {
			threat := SecurityThreat{
				Type:        "unusual_access_time",
				Severity:    SeverityWarning,
				Description: fmt.Sprintf("User %d has %d activities during off-hours", userID, offHoursCount),
				UserID:      &userID,
				Count:       offHoursCount,
				FirstSeen:   activities[len(activities)-1].CreatedAt,
				LastSeen:    activities[0].CreatedAt,
				Details: map[string]interface{}{
					"off_hours_activities": offHoursCount,
					"total_activities":     len(activities),
				},
			}
			threats = append(threats, threat)
		}

		// Alert if access from many different IPs
		if len(uniqueIPs) > 5 {
			threat := SecurityThreat{
				Type:        "multiple_ip_access",
				Severity:    SeverityWarning,
				Description: fmt.Sprintf("User %d accessed from %d different IP addresses", userID, len(uniqueIPs)),
				UserID:      &userID,
				Count:       len(uniqueIPs),
				FirstSeen:   activities[len(activities)-1].CreatedAt,
				LastSeen:    activities[0].CreatedAt,
				Details: map[string]interface{}{
					"unique_ips":       len(uniqueIPs),
					"total_activities": len(activities),
				},
			}
			threats = append(threats, threat)
		}
	}

	return threats, nil
}

// detectPrivilegeEscalation detects potential privilege escalation attempts
func (s *AuditService) detectPrivilegeEscalation() ([]SecurityThreat, error) {
	var threats []SecurityThreat

	// Get recent permission denied events
	since := time.Now().Add(-1 * time.Hour)
	logs, _, err := s.store.GetAuditLogsByAction("permission_denied", 0, 500)
	if err != nil {
		return threats, err
	}

	// Group by user
	userDenials := make(map[uint][]store.AuditLog)
	for _, log := range logs {
		if log.CreatedAt.After(since) && log.UserID != nil {
			userDenials[*log.UserID] = append(userDenials[*log.UserID], *log)
		}
	}

	// Check for excessive permission denials
	for userID, denials := range userDenials {
		if len(denials) >= 20 { // 20 or more denials in an hour
			threat := SecurityThreat{
				Type:        "privilege_escalation_attempt",
				Severity:    SeverityError,
				Description: fmt.Sprintf("User %d has %d permission denials in the last hour", userID, len(denials)),
				UserID:      &userID,
				Count:       len(denials),
				FirstSeen:   denials[len(denials)-1].CreatedAt,
				LastSeen:    denials[0].CreatedAt,
				Details: map[string]interface{}{
					"permission_denials": len(denials),
					"time_window":        "1 hour",
				},
			}
			threats = append(threats, threat)
		}
	}

	return threats, nil
}

// detectBruteForceAttacks detects brute force attacks across different attack vectors
func (s *AuditService) detectBruteForceAttacks() ([]SecurityThreat, error) {
	var threats []SecurityThreat

	// Get recent audit logs
	since := time.Now().Add(-30 * time.Minute)
	logs, _, err := s.store.ListAuditLogs(0, 1000)
	if err != nil {
		return threats, err
	}

	// Group by IP and count different types of failures
	ipActivity := make(map[string]map[string]int)
	ipFirstSeen := make(map[string]time.Time)
	ipLastSeen := make(map[string]time.Time)

	for _, log := range logs {
		if log.CreatedAt.After(since) && log.IPAddress != "" {
			if ipActivity[log.IPAddress] == nil {
				ipActivity[log.IPAddress] = make(map[string]int)
			}

			// Count failed actions
			if strings.Contains(log.Action, "failed") || strings.Contains(log.Action, "denied") {
				ipActivity[log.IPAddress][log.Action]++
			}

			// Track time range
			if ipFirstSeen[log.IPAddress].IsZero() || log.CreatedAt.Before(ipFirstSeen[log.IPAddress]) {
				ipFirstSeen[log.IPAddress] = log.CreatedAt
			}
			if ipLastSeen[log.IPAddress].IsZero() || log.CreatedAt.After(ipLastSeen[log.IPAddress]) {
				ipLastSeen[log.IPAddress] = log.CreatedAt
			}
		}
	}

	// Analyze patterns
	for ip, activities := range ipActivity {
		totalFailures := 0
		for _, count := range activities {
			totalFailures += count
		}

		if totalFailures >= 15 { // 15 or more failures in 30 minutes
			threat := SecurityThreat{
				Type:        "brute_force_attack",
				Severity:    SeverityCritical,
				Description: fmt.Sprintf("Potential brute force attack from IP %s with %d failed attempts", ip, totalFailures),
				IPAddress:   ip,
				Count:       totalFailures,
				FirstSeen:   ipFirstSeen[ip],
				LastSeen:    ipLastSeen[ip],
				Details: map[string]interface{}{
					"total_failures": totalFailures,
					"attack_types":   activities,
					"time_window":    "30 minutes",
				},
			}
			threats = append(threats, threat)
		}
	}

	return threats, nil
}

// GetAuditReport generates an audit report for a specific time period
func (s *AuditService) GetAuditReport(startTime, endTime time.Time, userID *uint) (*AuditReport, error) {
	report := &AuditReport{
		StartTime: startTime,
		EndTime:   endTime,
		UserID:    userID,
	}

	// Get audit logs for the period
	var logs []*store.AuditLog
	var err error

	if userID != nil {
		logs, _, err = s.store.GetAuditLogsByUserID(*userID, 0, 10000)
	} else {
		logs, _, err = s.store.ListAuditLogs(0, 10000)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Filter by time range and analyze
	filteredLogs := make([]*store.AuditLog, 0)
	actionCounts := make(map[string]int)
	userCounts := make(map[uint]int)
	ipCounts := make(map[string]int)

	for _, log := range logs {
		if log.CreatedAt.After(startTime) && log.CreatedAt.Before(endTime) {
			filteredLogs = append(filteredLogs, log)
			actionCounts[log.Action]++
			if log.UserID != nil {
				userCounts[*log.UserID]++
			}
			if log.IPAddress != "" {
				ipCounts[log.IPAddress]++
			}
		}
	}

	report.TotalEvents = len(filteredLogs)
	report.Events = filteredLogs
	report.ActionSummary = actionCounts
	report.UserActivity = userCounts
	report.IPActivity = ipCounts

	// Calculate statistics
	report.LoginAttempts = actionCounts["login"] + actionCounts["login_failed"]
	report.FailedLogins = actionCounts["login_failed"]
	report.PermissionDenials = actionCounts["permission_denied"]

	if report.LoginAttempts > 0 {
		report.LoginSuccessRate = float64(actionCounts["login"]) / float64(report.LoginAttempts) * 100
	}

	return report, nil
}

// AuditReport represents an audit report
type AuditReport struct {
	StartTime         time.Time         `json:"start_time"`
	EndTime           time.Time         `json:"end_time"`
	UserID            *uint             `json:"user_id,omitempty"`
	TotalEvents       int               `json:"total_events"`
	LoginAttempts     int               `json:"login_attempts"`
	FailedLogins      int               `json:"failed_logins"`
	LoginSuccessRate  float64           `json:"login_success_rate"`
	PermissionDenials int               `json:"permission_denials"`
	Events            []*store.AuditLog `json:"events"`
	ActionSummary     map[string]int    `json:"action_summary"`
	UserActivity      map[uint]int      `json:"user_activity"`
	IPActivity        map[string]int    `json:"ip_activity"`
}

// GetSecurityMetrics returns security metrics for monitoring
func (s *AuditService) GetSecurityMetrics(period time.Duration) (*SecurityMetrics, error) {
	since := time.Now().Add(-period)

	logs, _, err := s.store.ListAuditLogs(0, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	metrics := &SecurityMetrics{
		Period:    period,
		Timestamp: time.Now(),
	}

	// Analyze logs
	for _, log := range logs {
		if log.CreatedAt.After(since) {
			metrics.TotalEvents++

			switch log.Action {
			case "login":
				metrics.SuccessfulLogins++
			case "login_failed":
				metrics.FailedLogins++
			case "permission_denied":
				metrics.PermissionDenials++
			}

			if strings.Contains(log.Action, "failed") || strings.Contains(log.Action, "denied") {
				metrics.SecurityViolations++
			}
		}
	}

	// Calculate rates
	hours := period.Hours()
	if hours > 0 {
		metrics.EventsPerHour = float64(metrics.TotalEvents) / hours
		metrics.FailedLoginsPerHour = float64(metrics.FailedLogins) / hours
	}

	return metrics, nil
}

// SecurityMetrics represents security metrics
type SecurityMetrics struct {
	Period              time.Duration `json:"period"`
	Timestamp           time.Time     `json:"timestamp"`
	TotalEvents         int           `json:"total_events"`
	SuccessfulLogins    int           `json:"successful_logins"`
	FailedLogins        int           `json:"failed_logins"`
	PermissionDenials   int           `json:"permission_denials"`
	SecurityViolations  int           `json:"security_violations"`
	EventsPerHour       float64       `json:"events_per_hour"`
	FailedLoginsPerHour float64       `json:"failed_logins_per_hour"`
}

// StartMonitoring starts the continuous monitoring process
func (s *AuditService) StartMonitoring() {
	// Run anomaly detection every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			threats, err := s.DetectAnomalousActivity()
			if err != nil {
				fmt.Printf("Error detecting anomalous activity: %v\n", err)
				continue
			}

			// Log detected threats
			for _, threat := range threats {
				s.LogSecurityEvent(SecurityEvent{
					Type:      threat.Type,
					Severity:  string(threat.Severity),
					IPAddress: threat.IPAddress,
					UserID:    threat.UserID,
					Username:  threat.Username,
					Resource:  "security_monitoring",
					Action:    "threat_detected",
					Result:    "detected",
					Details: map[string]interface{}{
						"threat_description": threat.Description,
						"threat_count":       threat.Count,
						"first_seen":         threat.FirstSeen,
						"last_seen":          threat.LastSeen,
						"threat_details":     threat.Details,
					},
				})
			}

			// In a real implementation, you might want to:
			// - Send alerts to administrators
			// - Trigger automated responses
			// - Update security dashboards
		}
	}()
}

// Helper methods for audit handler

// GetAuditLogsByUserID gets audit logs for a specific user with pagination
func (s *AuditService) GetAuditLogsByUserID(userID uint, offset, limit int) (interface{}, int64, error) {
	return s.store.GetAuditLogsByUserID(userID, offset, limit)
}

// GetAuditLogsByAction gets audit logs for a specific action with pagination
func (s *AuditService) GetAuditLogsByAction(action string, offset, limit int) (interface{}, int64, error) {
	return s.store.GetAuditLogsByAction(action, offset, limit)
}

// GetAllAuditLogs gets all audit logs with pagination
func (s *AuditService) GetAllAuditLogs(offset, limit int) (interface{}, int64, error) {
	return s.store.ListAuditLogs(offset, limit)
}
