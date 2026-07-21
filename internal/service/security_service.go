package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/google/uuid"
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

// CreateSession creates a new user session persisted in the store
func (s *SecurityService) CreateSession(userID uint, ipAddress, userAgent string) (string, error) {
	sessionID := generateSessionID()
	now := time.Now()
	absoluteTimeout := s.config.Security.Session.AbsoluteTimeout
	if absoluteTimeout <= 0 {
		absoluteTimeout = 8 * time.Hour
	}

	// Enforce concurrent session limit
	if s.config.Security.Session.MaxConcurrentSessions > 0 {
		existing, err := s.store.GetUserSessions(userID)
		if err == nil && len(existing) >= s.config.Security.Session.MaxConcurrentSessions {
			overflow := len(existing) - s.config.Security.Session.MaxConcurrentSessions + 1
			// Evict oldest sessions first
			for overflow > 0 && len(existing) > 0 {
				oldestIdx := 0
				for i := 1; i < len(existing); i++ {
					if existing[i].CreatedAt.Before(existing[oldestIdx].CreatedAt) {
						oldestIdx = i
					}
				}
				_ = s.InvalidateSession(existing[oldestIdx].SessionID)
				existing = append(existing[:oldestIdx], existing[oldestIdx+1:]...)
				overflow--
			}
		}
	}

	userSession := &store.UserSession{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: now,
		LastSeen:  now,
		ExpiresAt: now.Add(absoluteTimeout),
		IsActive:  true,
	}
	if err := s.store.CreateUserSession(userSession); err != nil {
		return "", fmt.Errorf("failed to persist session: %w", err)
	}

	auditLog := &store.AuditLog{
		UserID:     &userID,
		Action:     "session_created",
		Resource:   "session",
		ResourceID: sessionID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    "New session created",
	}
	_ = s.store.CreateAuditLog(auditLog)

	return sessionID, nil
}

// ValidateSession validates and updates session activity
func (s *SecurityService) ValidateSession(sessionID string) (*SessionInfo, error) {
	userSession, err := s.store.GetUserSession(sessionID)
	if err != nil || userSession == nil {
		return nil, errors.New("session not found")
	}
	if !userSession.IsActive {
		return nil, errors.New("session not found")
	}

	now := time.Now()
	if now.After(userSession.ExpiresAt) {
		_ = s.InvalidateSession(sessionID)
		return nil, errors.New("session has expired")
	}

	if s.config.Security.Session.IdleTimeout > 0 {
		idleExpiry := userSession.LastSeen.Add(s.config.Security.Session.IdleTimeout)
		if now.After(idleExpiry) {
			_ = s.InvalidateSession(sessionID)
			return nil, errors.New("session has been idle too long")
		}
	}

	userSession.LastSeen = now
	if err := s.store.UpdateUserSession(userSession); err != nil {
		return nil, fmt.Errorf("failed to update session activity: %w", err)
	}

	return toSessionInfo(userSession), nil
}

// ValidateSessionID implements auth.SessionValidator
func (s *SecurityService) ValidateSessionID(sessionID string) error {
	_, err := s.ValidateSession(sessionID)
	return err
}

// InvalidateSession removes a session
func (s *SecurityService) InvalidateSession(sessionID string) error {
	userSession, err := s.store.GetUserSession(sessionID)
	if err != nil || userSession == nil {
		return nil
	}

	if err := s.store.DeleteUserSession(sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	auditLog := &store.AuditLog{
		UserID:     &userSession.UserID,
		Action:     "session_invalidated",
		Resource:   "session",
		ResourceID: sessionID,
		IPAddress:  userSession.IPAddress,
		UserAgent:  userSession.UserAgent,
		Details:    "Session invalidated",
	}
	_ = s.store.CreateAuditLog(auditLog)

	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a user
func (s *SecurityService) InvalidateAllUserSessions(userID uint) error {
	sessions, err := s.store.GetUserSessions(userID)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		_ = s.InvalidateSession(session.SessionID)
	}
	// Fallback bulk delete in case any remained
	return s.store.DeleteUserSessions(userID)
}

// GetUserSessions returns all active sessions for a user
func (s *SecurityService) GetUserSessions(userID uint) []*SessionInfo {
	sessions, err := s.store.GetUserSessions(userID)
	if err != nil {
		return []*SessionInfo{}
	}

	result := make([]*SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, toSessionInfo(session))
	}
	return result
}

// CleanupExpiredSessions removes expired sessions (should be called periodically)
func (s *SecurityService) CleanupExpiredSessions() {
	now := time.Now()
	_ = s.store.CleanupExpiredSessions(now)

	// Also expire idle sessions that are still marked active
	// Best-effort: scan active sessions per known users is expensive; rely on ValidateSession for idle checks.
}

func toSessionInfo(session *store.UserSession) *SessionInfo {
	return &SessionInfo{
		UserID:    session.UserID,
		SessionID: session.SessionID,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		CreatedAt: session.CreatedAt,
		LastSeen:  session.LastSeen,
		ExpiresAt: session.ExpiresAt,
	}
}

func generateSessionID() string {
	return "sess_" + uuid.NewString()
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
