package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

// SecurityService provides security-related functionality
type SecurityService struct {
	store  store.Store
	config *configs.Config
}

// NewSecurityService creates a new SecurityService instance
func NewSecurityService(store store.Store, config *configs.Config) *SecurityService {
	return &SecurityService{
		store:  store,
		config: config,
	}
}

// PasswordValidationError represents password validation errors
type PasswordValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e PasswordValidationError) Error() string {
	return e.Message
}

// ValidatePassword validates password against security policy
func (s *SecurityService) ValidatePassword(password string) []PasswordValidationError {
	var errors []PasswordValidationError
	policy := s.config.Security.Password

	// Check minimum length
	if len(password) < policy.MinLength {
		errors = append(errors, PasswordValidationError{
			Field:   "password",
			Message: fmt.Sprintf("Password must be at least %d characters long", policy.MinLength),
		})
	}

	// Check uppercase requirement
	if policy.RequireUppercase {
		if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
			errors = append(errors, PasswordValidationError{
				Field:   "password",
				Message: "Password must contain at least one uppercase letter",
			})
		}
	}

	// Check lowercase requirement
	if policy.RequireLowercase {
		if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
			errors = append(errors, PasswordValidationError{
				Field:   "password",
				Message: "Password must contain at least one lowercase letter",
			})
		}
	}

	// Check numbers requirement
	if policy.RequireNumbers {
		if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
			errors = append(errors, PasswordValidationError{
				Field:   "password",
				Message: "Password must contain at least one number",
			})
		}
	}

	// Check symbols requirement
	if policy.RequireSymbols {
		if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~`+"`"+`]`, password); !matched {
			errors = append(errors, PasswordValidationError{
				Field:   "password",
				Message: "Password must contain at least one special character",
			})
		}
	}

	// Check for common weak passwords
	weakPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "root", "user", "guest",
	}

	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if lowerPassword == weak {
			errors = append(errors, PasswordValidationError{
				Field:   "password",
				Message: "Password is too common and easily guessable",
			})
			break
		}
	}

	return errors
}

// CheckAccountLockout checks if an account is locked due to failed login attempts
func (s *SecurityService) CheckAccountLockout(userID uint) (bool, time.Time, error) {
	if !s.config.Security.AccountLock.Enabled {
		return false, time.Time{}, nil
	}

	// Get recent failed login attempts
	since := time.Now().Add(-s.config.Security.AccountLock.ResetAfter)
	attempts, _, err := s.store.GetAuditLogsByUserID(userID, 0, 100)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Count failed login attempts within the reset window
	failedCount := 0
	var lastFailedAttempt time.Time

	for _, attempt := range attempts {
		if attempt.CreatedAt.Before(since) {
			break // Older than reset window
		}

		if attempt.Action == "login_failed" {
			failedCount++
			if lastFailedAttempt.IsZero() || attempt.CreatedAt.After(lastFailedAttempt) {
				lastFailedAttempt = attempt.CreatedAt
			}
		} else if attempt.Action == "login" {
			// Successful login resets the counter
			break
		}
	}

	// Check if account should be locked
	if failedCount >= s.config.Security.AccountLock.MaxFailedAttempts {
		lockoutEnd := lastFailedAttempt.Add(s.config.Security.AccountLock.LockoutDuration)
		if time.Now().Before(lockoutEnd) {
			return true, lockoutEnd, nil
		}
	}

	return false, time.Time{}, nil
}

// RecordFailedLogin records a failed login attempt
func (s *SecurityService) RecordFailedLogin(userID *uint, username, ipAddress, userAgent string) error {
	auditLog := &store.AuditLog{
		UserID:     userID,
		Action:     "login_failed",
		Resource:   "user",
		ResourceID: username,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    fmt.Sprintf("Failed login attempt for username: %s", username),
	}

	return s.store.CreateAuditLog(auditLog)
}

// RecordSuccessfulLogin records a successful login
func (s *SecurityService) RecordSuccessfulLogin(userID uint, ipAddress, userAgent string) error {
	auditLog := &store.AuditLog{
		UserID:     &userID,
		Action:     "login",
		Resource:   "user",
		ResourceID: fmt.Sprintf("%d", userID),
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    "Successful login",
	}

	return s.store.CreateAuditLog(auditLog)
}

// SessionInfo represents active session information
type SessionInfo struct {
	UserID    uint      `json:"user_id"`
	SessionID string    `json:"session_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
	ExpiresAt time.Time `json:"expires_at"`
}

// In-memory session store (in production, this should be Redis or database)
var activeSessions = make(map[string]*SessionInfo)
var userSessions = make(map[uint][]string) // userID -> sessionIDs

// CreateSession creates a new user session
func (s *SecurityService) CreateSession(userID uint, ipAddress, userAgent string) (string, error) {
	sessionID := generateSessionID()
	now := time.Now()

	session := &SessionInfo{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: now,
		LastSeen:  now,
		ExpiresAt: now.Add(s.config.Security.Session.AbsoluteTimeout),
	}

	// Check concurrent session limit
	if s.config.Security.Session.MaxConcurrentSessions > 0 {
		userSessionIDs := userSessions[userID]
		if len(userSessionIDs) >= s.config.Security.Session.MaxConcurrentSessions {
			// Remove oldest session
			oldestSessionID := userSessionIDs[0]
			s.InvalidateSession(oldestSessionID)
		}
	}

	// Store session
	activeSessions[sessionID] = session
	userSessions[userID] = append(userSessions[userID], sessionID)

	// Record session creation
	auditLog := &store.AuditLog{
		UserID:     &userID,
		Action:     "session_created",
		Resource:   "session",
		ResourceID: sessionID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    "New session created",
	}
	s.store.CreateAuditLog(auditLog)

	return sessionID, nil
}

// ValidateSession validates and updates session activity
func (s *SecurityService) ValidateSession(sessionID string) (*SessionInfo, error) {
	session, exists := activeSessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	now := time.Now()

	// Check if session has expired (absolute timeout)
	if now.After(session.ExpiresAt) {
		s.InvalidateSession(sessionID)
		return nil, errors.New("session has expired")
	}

	// Check idle timeout
	if s.config.Security.Session.IdleTimeout > 0 {
		idleExpiry := session.LastSeen.Add(s.config.Security.Session.IdleTimeout)
		if now.After(idleExpiry) {
			s.InvalidateSession(sessionID)
			return nil, errors.New("session has been idle too long")
		}
	}

	// Update last seen time
	session.LastSeen = now

	return session, nil
}

// InvalidateSession removes a session
func (s *SecurityService) InvalidateSession(sessionID string) error {
	session, exists := activeSessions[sessionID]
	if !exists {
		return nil // Already invalid
	}

	// Remove from active sessions
	delete(activeSessions, sessionID)

	// Remove from user sessions
	userSessionIDs := userSessions[session.UserID]
	for i, id := range userSessionIDs {
		if id == sessionID {
			userSessions[session.UserID] = append(userSessionIDs[:i], userSessionIDs[i+1:]...)
			break
		}
	}

	// Clean up empty user session list
	if len(userSessions[session.UserID]) == 0 {
		delete(userSessions, session.UserID)
	}

	// Record session invalidation
	auditLog := &store.AuditLog{
		UserID:     &session.UserID,
		Action:     "session_invalidated",
		Resource:   "session",
		ResourceID: sessionID,
		IPAddress:  session.IPAddress,
		UserAgent:  session.UserAgent,
		Details:    "Session invalidated",
	}
	s.store.CreateAuditLog(auditLog)

	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a user
func (s *SecurityService) InvalidateAllUserSessions(userID uint) error {
	sessionIDs := userSessions[userID]
	for _, sessionID := range sessionIDs {
		s.InvalidateSession(sessionID)
	}
	return nil
}

// GetUserSessions returns all active sessions for a user
func (s *SecurityService) GetUserSessions(userID uint) []*SessionInfo {
	sessionIDs := userSessions[userID]
	sessions := make([]*SessionInfo, 0, len(sessionIDs))

	for _, sessionID := range sessionIDs {
		if session, exists := activeSessions[sessionID]; exists {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// CleanupExpiredSessions removes expired sessions (should be called periodically)
func (s *SecurityService) CleanupExpiredSessions() {
	now := time.Now()

	for sessionID, session := range activeSessions {
		expired := false

		// Check absolute timeout
		if now.After(session.ExpiresAt) {
			expired = true
		}

		// Check idle timeout
		if s.config.Security.Session.IdleTimeout > 0 {
			idleExpiry := session.LastSeen.Add(s.config.Security.Session.IdleTimeout)
			if now.After(idleExpiry) {
				expired = true
			}
		}

		if expired {
			s.InvalidateSession(sessionID)
		}
	}
}

// Helper function to generate session ID
func generateSessionID() string {
	// In production, use a cryptographically secure random generator
	return fmt.Sprintf("sess_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}

// RecordSecurityEvent records a security-related event
func (s *SecurityService) RecordSecurityEvent(userID *uint, action, resource, resourceID, ipAddress, userAgent, details string) error {
	auditLog := &store.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
	}

	return s.store.CreateAuditLog(auditLog)
}

// GetSecurityEvents returns security events for analysis
func (s *SecurityService) GetSecurityEvents(page, pageSize int) ([]*store.AuditLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.store.ListAuditLogs(offset, pageSize)
}

// DetectSuspiciousActivity analyzes recent activities for suspicious patterns
func (s *SecurityService) DetectSuspiciousActivity(userID uint) ([]string, error) {
	var warnings []string

	// Get recent audit logs for the user
	logs, _, err := s.store.GetAuditLogsByUserID(userID, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Analyze patterns
	recentFailures := 0
	uniqueIPs := make(map[string]bool)
	now := time.Now()

	for _, log := range logs {
		// Only analyze recent events (last 24 hours)
		if log.CreatedAt.Before(now.Add(-24 * time.Hour)) {
			break
		}

		if log.Action == "login_failed" {
			recentFailures++
		}

		if log.IPAddress != "" {
			uniqueIPs[log.IPAddress] = true
		}
	}

	// Check for multiple failed logins
	if recentFailures > 3 {
		warnings = append(warnings, fmt.Sprintf("Multiple failed login attempts (%d) in the last 24 hours", recentFailures))
	}

	// Check for logins from multiple IPs
	if len(uniqueIPs) > 3 {
		warnings = append(warnings, fmt.Sprintf("Logins from multiple IP addresses (%d) in the last 24 hours", len(uniqueIPs)))
	}

	return warnings, nil
}
