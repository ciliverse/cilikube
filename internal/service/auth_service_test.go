package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestAuthService() (*AuthService, store.Store) {
	// Create test configuration
	config := &configs.Config{
		Security: configs.SecurityConfig{
			Password: configs.PasswordConfig{
				MinLength:        8,
				RequireUppercase: false,
				RequireLowercase: true,
				RequireNumbers:   true,
				RequireSymbols:   false,
			},
			AccountLock: configs.AccountLockConfig{
				Enabled:           true,
				MaxFailedAttempts: 5,
				LockoutDuration:   15 * time.Minute,
				ResetAfter:        1 * time.Hour,
			},
		},
	}

	// Create in-memory store
	testStore := store.NewMemoryStore()
	testStore.Initialize()

	// Create auth service
	authService := NewAuthService(testStore, config)

	return authService, testStore
}

func TestAuthService_Register(t *testing.T) {
	authService, _ := setupTestAuthService()

	tests := []struct {
		name        string
		request     *models.RegisterRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid registration",
			request: &models.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Weak password",
			request: &models.RegisterRequest{
				Username: "testuser2",
				Email:    "test2@example.com",
				Password: "weak",
			},
			expectError: true,
			errorMsg:    "password validation failed",
		},
		{
			name: "Duplicate username",
			request: &models.RegisterRequest{
				Username: "testuser",
				Email:    "different@example.com",
				Password: "password123",
			},
			expectError: true,
			errorMsg:    "username already exists",
		},
		{
			name: "Duplicate email",
			request: &models.RegisterRequest{
				Username: "differentuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectError: true,
			errorMsg:    "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Register(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.request.Username, response.Username)
				assert.Equal(t, tt.request.Email, response.Email)
				assert.True(t, response.IsActive)
				assert.False(t, response.EmailVerified)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test user
	testUser := &store.User{
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "password123",
		DisplayName:   "Test User",
		IsActive:      true,
		EmailVerified: false,
	}
	err := testStore.CreateUser(testUser)
	require.NoError(t, err)

	// Assign viewer role
	viewerRole, err := testStore.GetRoleByName("viewer")
	require.NoError(t, err)
	err = testStore.AssignRole(testUser.ID, viewerRole.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		request     *models.LoginRequest
		ipAddress   string
		userAgent   string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid login",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			ipAddress:   "127.0.0.1",
			userAgent:   "test-agent",
			expectError: false,
		},
		{
			name: "Invalid username",
			request: &models.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			ipAddress:   "127.0.0.1",
			userAgent:   "test-agent",
			expectError: true,
			errorMsg:    "invalid username or password",
		},
		{
			name: "Invalid password",
			request: &models.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			ipAddress:   "127.0.0.1",
			userAgent:   "test-agent",
			expectError: true,
			errorMsg:    "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.Login(tt.request, tt.ipAddress, tt.userAgent)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, "testuser", response.User.Username)
				assert.Equal(t, "viewer", response.User.Role)
				assert.True(t, response.ExpiresAt.After(time.Now()))
			}
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test user
	testUser := &store.User{
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "oldpassword123",
		DisplayName:   "Test User",
		IsActive:      true,
		EmailVerified: false,
	}
	err := testStore.CreateUser(testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		request     *models.ChangePasswordRequest
		expectError bool
		errorMsg    string
	}{
		{
			name:   "Valid password change",
			userID: testUser.ID,
			request: &models.ChangePasswordRequest{
				OldPassword: "oldpassword123",
				NewPassword: "newpassword123",
			},
			expectError: false,
		},
		{
			name:   "Wrong old password",
			userID: testUser.ID,
			request: &models.ChangePasswordRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword123",
			},
			expectError: true,
			errorMsg:    "old password is incorrect",
		},
		{
			name:   "Weak new password",
			userID: testUser.ID,
			request: &models.ChangePasswordRequest{
				OldPassword: "oldpassword123",
				NewPassword: "weak",
			},
			expectError: true,
			errorMsg:    "password validation failed",
		},
		{
			name:   "Non-existent user",
			userID: 999,
			request: &models.ChangePasswordRequest{
				OldPassword: "oldpassword123",
				NewPassword: "newpassword123",
			},
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.ChangePassword(tt.userID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify password was actually changed
				updatedUser, err := testStore.GetUserByID(tt.userID)
				assert.NoError(t, err)
				assert.True(t, updatedUser.CheckPassword(tt.request.NewPassword))
				assert.False(t, updatedUser.CheckPassword(tt.request.OldPassword))
			}
		})
	}
}

func TestAuthService_GetProfile(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test user
	testUser := &store.User{
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "password123",
		DisplayName:   "Test User",
		AvatarURL:     "https://example.com/avatar.jpg",
		IsActive:      true,
		EmailVerified: true,
	}
	err := testStore.CreateUser(testUser)
	require.NoError(t, err)

	// Assign admin role
	adminRole, err := testStore.GetRoleByName("admin")
	require.NoError(t, err)
	err = testStore.AssignRole(testUser.ID, adminRole.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid profile retrieval",
			userID:      testUser.ID,
			expectError: false,
		},
		{
			name:        "Non-existent user",
			userID:      999,
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := authService.GetProfile(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, "testuser", profile.Username)
				assert.Equal(t, "test@example.com", profile.Email)
				assert.Equal(t, "Test User", profile.DisplayName)
				assert.Equal(t, "https://example.com/avatar.jpg", profile.AvatarURL)
				assert.Equal(t, "admin", profile.Role)
				assert.True(t, profile.IsActive)
				assert.True(t, profile.EmailVerified)
			}
		})
	}
}

func TestAuthService_UpdateProfile(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test user
	testUser := &store.User{
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "password123",
		DisplayName:   "Test User",
		IsActive:      true,
		EmailVerified: false,
	}
	err := testStore.CreateUser(testUser)
	require.NoError(t, err)

	// Create another user to test email conflict
	anotherUser := &store.User{
		Username:      "anotheruser",
		Email:         "another@example.com",
		PasswordHash:  "password123",
		DisplayName:   "Another User",
		IsActive:      true,
		EmailVerified: false,
	}
	err = testStore.CreateUser(anotherUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		request     *models.UpdateProfileRequest
		expectError bool
		errorMsg    string
	}{
		{
			name:   "Valid profile update",
			userID: testUser.ID,
			request: &models.UpdateProfileRequest{
				Email:       "updated@example.com",
				DisplayName: "Updated User",
				AvatarURL:   "https://example.com/new-avatar.jpg",
			},
			expectError: false,
		},
		{
			name:   "Email conflict",
			userID: testUser.ID,
			request: &models.UpdateProfileRequest{
				Email:       "another@example.com",
				DisplayName: "Updated User",
				AvatarURL:   "",
			},
			expectError: true,
			errorMsg:    "email is already being used by another user",
		},
		{
			name:   "Non-existent user",
			userID: 999,
			request: &models.UpdateProfileRequest{
				Email:       "test@example.com",
				DisplayName: "Test",
				AvatarURL:   "",
			},
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.UpdateProfile(tt.userID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.request.Email, response.Email)
				assert.Equal(t, tt.request.DisplayName, response.DisplayName)
				assert.Equal(t, tt.request.AvatarURL, response.AvatarURL)

				// Verify changes were persisted
				updatedUser, err := testStore.GetUserByID(tt.userID)
				assert.NoError(t, err)
				assert.Equal(t, tt.request.Email, updatedUser.Email)
				assert.Equal(t, tt.request.DisplayName, updatedUser.DisplayName)
				assert.Equal(t, tt.request.AvatarURL, updatedUser.AvatarURL)
			}
		})
	}
}

func TestAuthService_AccountLockout(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test user
	testUser := &store.User{
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "password123",
		DisplayName:   "Test User",
		IsActive:      true,
		EmailVerified: false,
	}
	err := testStore.CreateUser(testUser)
	require.NoError(t, err)

	// Assign viewer role
	viewerRole, err := testStore.GetRoleByName("viewer")
	require.NoError(t, err)
	err = testStore.AssignRole(testUser.ID, viewerRole.ID)
	require.NoError(t, err)

	// Simulate multiple failed login attempts
	loginRequest := &models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// First 4 attempts should fail normally
	for i := 0; i < 4; i++ {
		_, err := authService.Login(loginRequest, "127.0.0.1", "test-agent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid username or password")
	}

	// 5th attempt should trigger account lockout
	_, err = authService.Login(loginRequest, "127.0.0.1", "test-agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")

	// 6th attempt should show account is locked
	_, err = authService.Login(loginRequest, "127.0.0.1", "test-agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account is temporarily locked")

	// Even correct password should be rejected when locked
	correctRequest := &models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	_, err = authService.Login(correctRequest, "127.0.0.1", "test-agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account is temporarily locked")
}

func TestAuthService_ValidatePassword(t *testing.T) {
	authService, _ := setupTestAuthService()

	tests := []struct {
		name           string
		password       string
		expectValid    bool
		expectedErrors int
	}{
		{
			name:           "Valid password",
			password:       "password123",
			expectValid:    true,
			expectedErrors: 0,
		},
		{
			name:           "Too short",
			password:       "pass1",
			expectValid:    false,
			expectedErrors: 1,
		},
		{
			name:           "No numbers",
			password:       "password",
			expectValid:    false,
			expectedErrors: 1,
		},
		{
			name:           "Common weak password",
			password:       "password",
			expectValid:    false,
			expectedErrors: 2, // No numbers + common password
		},
		{
			name:           "Multiple issues",
			password:       "123",
			expectValid:    false,
			expectedErrors: 2, // Too short + no lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := authService.ValidatePassword(tt.password)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tt.expectValid, response.Valid)
			assert.Equal(t, tt.expectedErrors, len(response.Errors))
		})
	}
}

func TestAuthService_UserManagement(t *testing.T) {
	authService, testStore := setupTestAuthService()

	// Create test users
	for i := 1; i <= 5; i++ {
		testUser := &store.User{
			Username:      fmt.Sprintf("user%d", i),
			Email:         fmt.Sprintf("user%d@example.com", i),
			PasswordHash:  "password123",
			DisplayName:   fmt.Sprintf("User %d", i),
			IsActive:      i%2 == 1, // Alternate active/inactive
			EmailVerified: false,
		}
		err := testStore.CreateUser(testUser)
		require.NoError(t, err)
	}

	t.Run("Get user list", func(t *testing.T) {
		users, total, err := authService.GetUserList(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total) // 5 test users + 1 admin from initialization
		assert.Len(t, users, 6)
	})

	t.Run("Get user list with pagination", func(t *testing.T) {
		users, total, err := authService.GetUserList(1, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total)
		assert.Len(t, users, 3)

		users, total, err = authService.GetUserList(2, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total)
		assert.Len(t, users, 3)
	})

	t.Run("Update user status", func(t *testing.T) {
		// Get a user to update
		users, _, err := authService.GetUserList(1, 1)
		require.NoError(t, err)
		require.Len(t, users, 1)

		userID := users[0].ID
		originalStatus := users[0].IsActive

		// Update status
		err = authService.UpdateUserStatus(userID, !originalStatus)
		assert.NoError(t, err)

		// Verify status was changed
		updatedUser, err := testStore.GetUserByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, !originalStatus, updatedUser.IsActive)
	})

	t.Run("Delete user", func(t *testing.T) {
		// Get a user to delete
		users, _, err := authService.GetUserList(1, 1)
		require.NoError(t, err)
		require.Len(t, users, 1)

		userID := users[0].ID

		// Delete user
		err = authService.DeleteUser(userID)
		assert.NoError(t, err)

		// Verify user was deleted
		_, err = testStore.GetUserByID(userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
