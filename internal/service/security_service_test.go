package service

import (
	"testing"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

func TestPasswordValidation(t *testing.T) {
	// Create test config with security settings
	config := &configs.Config{
		Security: configs.SecurityConfig{
			Password: configs.PasswordConfig{
				MinLength:        8,
				RequireUppercase: true,
				RequireLowercase: true,
				RequireNumbers:   true,
				RequireSymbols:   false,
			},
		},
	}

	// Create memory store and security service
	memStore := store.NewMemoryStore()
	securityService := NewSecurityService(memStore, config)

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "Valid password",
			password: "Password123",
			valid:    true,
		},
		{
			name:     "Too short",
			password: "Pass1",
			valid:    false,
		},
		{
			name:     "No uppercase",
			password: "password123",
			valid:    false,
		},
		{
			name:     "No lowercase",
			password: "PASSWORD123",
			valid:    false,
		},
		{
			name:     "No numbers",
			password: "Password",
			valid:    false,
		},
		{
			name:     "Common weak password",
			password: "password",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := securityService.ValidatePassword(tt.password)
			isValid := len(errors) == 0

			if isValid != tt.valid {
				t.Errorf("ValidatePassword() = %v, want %v", isValid, tt.valid)
				if len(errors) > 0 {
					t.Logf("Validation errors: %+v", errors)
				}
			}
		})
	}
}

func TestSessionManagement(t *testing.T) {
	// Create test config
	config := &configs.Config{
		Security: configs.SecurityConfig{
			Session: configs.SessionConfig{
				MaxConcurrentSessions: 2,
				IdleTimeout:           30 * time.Minute,
				AbsoluteTimeout:       8 * time.Hour,
			},
		},
	}

	// Create memory store and security service
	memStore := store.NewMemoryStore()
	securityService := NewSecurityService(memStore, config)

	userID := uint(1)
	ipAddress := "192.168.1.1"
	userAgent := "Test Browser"

	// Test session creation
	sessionID1, err := securityService.CreateSession(userID, ipAddress, userAgent)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Test session validation
	session, err := securityService.ValidateSession(sessionID1)
	if err != nil {
		t.Fatalf("Failed to validate session: %v", err)
	}

	if session.UserID != userID {
		t.Errorf("Session user ID = %v, want %v", session.UserID, userID)
	}

	// Test session invalidation
	err = securityService.InvalidateSession(sessionID1)
	if err != nil {
		t.Fatalf("Failed to invalidate session: %v", err)
	}

	// Test that invalidated session is no longer valid
	_, err = securityService.ValidateSession(sessionID1)
	if err == nil {
		t.Error("Expected error when validating invalidated session")
	}
}

func TestAccountLockout(t *testing.T) {
	// Create test config with account lockout enabled
	config := &configs.Config{
		Security: configs.SecurityConfig{
			AccountLock: configs.AccountLockConfig{
				Enabled:           true,
				MaxFailedAttempts: 3,
				LockoutDuration:   15 * time.Minute,
				ResetAfter:        1 * time.Hour,
			},
		},
	}

	// Create memory store and security service
	memStore := store.NewMemoryStore()
	if err := memStore.Initialize(); err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	securityService := NewSecurityService(memStore, config)

	userID := uint(1)
	username := "testuser"
	ipAddress := "192.168.1.1"
	userAgent := "Test Browser"

	// Record multiple failed login attempts
	for i := 0; i < 3; i++ {
		err := securityService.RecordFailedLogin(&userID, username, ipAddress, userAgent)
		if err != nil {
			t.Fatalf("Failed to record failed login: %v", err)
		}
	}

	// Check if account is locked
	isLocked, lockoutEnd, err := securityService.CheckAccountLockout(userID)
	if err != nil {
		t.Fatalf("Failed to check account lockout: %v", err)
	}

	if !isLocked {
		t.Error("Expected account to be locked after 3 failed attempts")
	}

	if lockoutEnd.IsZero() {
		t.Error("Expected lockout end time to be set")
	}

	t.Logf("Account locked until: %v", lockoutEnd)
}

func TestSuspiciousActivityDetection(t *testing.T) {
	// Create test config
	config := &configs.Config{}

	// Create memory store and security service
	memStore := store.NewMemoryStore()
	if err := memStore.Initialize(); err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	securityService := NewSecurityService(memStore, config)

	userID := uint(1)

	// Create multiple failed login attempts from different IPs
	ips := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1", "203.0.113.1"}
	for _, ip := range ips {
		for i := 0; i < 2; i++ {
			err := securityService.RecordFailedLogin(&userID, "testuser", ip, "Test Browser")
			if err != nil {
				t.Fatalf("Failed to record failed login: %v", err)
			}
		}
	}

	// Detect suspicious activity
	warnings, err := securityService.DetectSuspiciousActivity(userID)
	if err != nil {
		t.Fatalf("Failed to detect suspicious activity: %v", err)
	}

	if len(warnings) == 0 {
		t.Error("Expected suspicious activity warnings")
	}

	t.Logf("Suspicious activity warnings: %v", warnings)
}
