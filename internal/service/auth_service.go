package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/auth"
)

// AuthService provides authentication and user management functionality
type AuthService struct {
	store           store.Store
	config          *configs.Config
	securityService *SecurityService
	auditService    *AuditService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(store store.Store, config *configs.Config) *AuthService {
	securityService := NewSecurityService(store, config)
	auditService := NewAuditService(store, config)
	return &AuthService{
		store:           store,
		config:          config,
		securityService: securityService,
		auditService:    auditService,
	}
}

// Login authenticates a user with username/password and returns JWT token
func (s *AuthService) Login(req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Get user by username
	storeUser, err := s.store.GetUserByUsername(req.Username)
	if err != nil {
		// Record failed login attempt for unknown user
		s.securityService.RecordFailedLogin(nil, req.Username, ipAddress, userAgent)
		s.auditService.LogAuthenticationEvent(AuditEventType("login_failed"), nil, req.Username, ipAddress, userAgent, false, map[string]interface{}{
			"reason": "user_not_found",
		})
		return nil, errors.New("invalid username or password")
	}

	// Check account lockout status
	isLocked, lockoutEnd, err := s.securityService.CheckAccountLockout(storeUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account lockout: %w", err)
	}
	if isLocked {
		return nil, fmt.Errorf("account is temporarily locked until %s due to multiple failed login attempts", lockoutEnd.Format("2006-01-02 15:04:05"))
	}

	// Check if user is active
	if !storeUser.IsActive {
		s.securityService.RecordFailedLogin(&storeUser.ID, req.Username, ipAddress, userAgent)
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !storeUser.CheckPassword(req.Password) {
		s.securityService.RecordFailedLogin(&storeUser.ID, req.Username, ipAddress, userAgent)
		return nil, errors.New("invalid username or password")
	}

	// Record successful login
	if err := s.securityService.RecordSuccessfulLogin(storeUser.ID, ipAddress, userAgent); err != nil {
		fmt.Printf("Failed to record successful login: %v\n", err)
	}

	// Update last login time
	now := time.Now()
	storeUser.LastLoginAt = &now
	if err := s.store.UpdateUser(storeUser); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Convert store user to models user for JWT generation
	user := s.convertStoreUserToModelsUser(storeUser)

	// Get user roles for JWT token
	roles, err := s.store.GetUserRoles(storeUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Set primary role (for backward compatibility)
	if len(roles) > 0 {
		user.Role = roles[0].Name
	} else {
		user.Role = "viewer" // Default role
	}

	// Generate JWT token
	token, expiresAt, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session if session management is enabled
	sessionID, err := s.securityService.CreateSession(storeUser.ID, ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
	}

	// Create audit log
	s.createAuditLog(&storeUser.ID, "login", "user", fmt.Sprintf("%d", storeUser.ID), ipAddress, userAgent, fmt.Sprintf("User logged in successfully, session: %s", sessionID))

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// RefreshToken generates a new JWT token from a valid existing token
func (s *AuthService) RefreshToken(tokenString string) (*models.TokenResponse, error) {
	// Parse the existing token
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if token is close to expiry (within 1 hour)
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return nil, errors.New("token is not eligible for refresh yet")
	}

	// Get user from store to ensure they still exist and are active
	storeUser, err := s.store.GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !storeUser.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Convert to models user
	user := s.convertStoreUserToModelsUser(storeUser)

	// Get current roles
	roles, err := s.store.GetUserRoles(storeUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) > 0 {
		user.Role = roles[0].Name
	} else {
		user.Role = "viewer"
	}

	// Generate new token
	newToken, expiresAt, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	// Create audit log
	s.createAuditLog(&storeUser.ID, "token_refresh", "user", fmt.Sprintf("%d", storeUser.ID), "", "", "Token refreshed successfully")

	return &models.TokenResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
	}, nil
}

// Logout invalidates a user session (placeholder for future session management)
func (s *AuthService) Logout(userID uint) error {
	// Create audit log
	s.createAuditLog(&userID, "logout", "user", fmt.Sprintf("%d", userID), "", "", "User logged out")

	// In the future, we could implement token blacklisting here
	return nil
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// Validate password against security policy
	if validationErrors := s.securityService.ValidatePassword(req.Password); len(validationErrors) > 0 {
		return nil, fmt.Errorf("password validation failed: %s", validationErrors[0].Message)
	}

	// Check if username already exists
	_, err := s.store.GetUserByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.store.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Create new store user
	storeUser := &store.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  req.Password, // Will be hashed by store
		DisplayName:   req.Username,
		IsActive:      true,
		EmailVerified: false,
	}

	// Create user in store
	if err := s.store.CreateUser(storeUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default viewer role
	viewerRole, err := s.store.GetRoleByName("viewer")
	if err != nil {
		return nil, fmt.Errorf("failed to get viewer role: %w", err)
	}

	if err := s.store.AssignRole(storeUser.ID, viewerRole.ID); err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	// Create audit log
	s.createAuditLog(nil, "user_register", "user", fmt.Sprintf("%d", storeUser.ID), "", "", "New user registered")

	// Convert to response
	user := s.convertStoreUserToModelsUser(storeUser)
	response := user.ToResponse()
	return &response, nil
}

// GetProfile gets detailed user profile information
func (s *AuthService) GetProfile(userID uint) (*models.UserProfileResponse, error) {
	// Get user from store
	storeUser, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Convert to models user
	user := s.convertStoreUserToModelsUser(storeUser)

	// Get user roles
	roles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Set primary role and roles array
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	if len(roles) > 0 {
		user.Role = roles[0].Name
	} else {
		user.Role = "viewer"
		roleNames = []string{"viewer"}
	}

	// Get OAuth providers
	oauthProviders, err := s.store.ListUserOAuthProviders(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth providers: %w", err)
	}

	// Convert to profile response
	profile := user.ToProfileResponse()

	// Set roles array
	profile.Roles = roleNames

	// Add OAuth provider info
	for _, provider := range oauthProviders {
		profile.OAuthProviders = append(profile.OAuthProviders, models.OAuthProviderInfo{
			Provider:       provider.Provider,
			ProviderUserID: provider.ProviderUserID,
			ConnectedAt:    provider.CreatedAt,
			ExpiresAt:      provider.ExpiresAt,
		})
	}

	return &profile, nil
}

// GetProfileLegacy returns profile in legacy format for backward compatibility
func (s *AuthService) GetProfileLegacy(userID uint) (map[string]interface{}, error) {
	profile, err := s.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	// Get user roles for legacy format
	roles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return map[string]interface{}{
		"username": profile.Username,
		"roles":    roleNames,
	}, nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(userID uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	// Get user from store
	storeUser, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and if it conflicts
	if req.Email != storeUser.Email {
		_, err := s.store.GetUserByEmail(req.Email)
		if err == nil {
			return nil, errors.New("email is already being used by another user")
		}
	}

	// Update user information
	storeUser.Email = req.Email
	storeUser.DisplayName = req.DisplayName
	storeUser.AvatarURL = req.AvatarURL

	if err := s.store.UpdateUser(storeUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	s.createAuditLog(&userID, "profile_update", "user", fmt.Sprintf("%d", userID), "", "", "User profile updated")

	// Convert and return response
	user := s.convertStoreUserToModelsUser(storeUser)
	response := user.ToResponse()
	return &response, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	// Get user from store
	storeUser, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !storeUser.CheckPassword(req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// Validate new password against security policy
	if validationErrors := s.securityService.ValidatePassword(req.NewPassword); len(validationErrors) > 0 {
		return fmt.Errorf("password validation failed: %s", validationErrors[0].Message)
	}

	// Update password
	if err := storeUser.HashPassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := s.store.UpdateUser(storeUser); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all user sessions to force re-login with new password
	if err := s.securityService.InvalidateAllUserSessions(userID); err != nil {
		fmt.Printf("Failed to invalidate user sessions: %v\n", err)
	}

	// Create audit log
	s.createAuditLog(&userID, "password_change", "user", fmt.Sprintf("%d", userID), "", "", "User password changed")

	return nil
}

// GetUserList gets paginated user list (admin function)
func (s *AuthService) GetUserList(page, pageSize int) ([]models.UserResponse, int64, error) {
	offset := (page - 1) * pageSize
	storeUsers, total, err := s.store.ListUsers(offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user list: %w", err)
	}

	// Convert to response format
	responses := make([]models.UserResponse, len(storeUsers))
	for i, storeUser := range storeUsers {
		user := s.convertStoreUserToModelsUser(storeUser)

		// Get user roles for primary role
		roles, err := s.store.GetUserRoles(storeUser.ID)
		if err == nil && len(roles) > 0 {
			user.Role = roles[0].Name
		} else {
			user.Role = "viewer"
		}

		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// UpdateUserStatus updates user active status (admin function)
func (s *AuthService) UpdateUserStatus(userID uint, isActive bool) error {
	storeUser, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	storeUser.IsActive = isActive
	if err := s.store.UpdateUser(storeUser); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Create audit log
	status := "activated"
	if !isActive {
		status = "deactivated"
	}
	s.createAuditLog(nil, "user_status_change", "user", fmt.Sprintf("%d", userID), "", "", fmt.Sprintf("User %s", status))

	return nil
}

// DeleteUser deletes a user (admin function)
func (s *AuthService) DeleteUser(userID uint) error {
	if err := s.store.DeleteUser(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Create audit log
	s.createAuditLog(nil, "user_delete", "user", fmt.Sprintf("%d", userID), "", "", "User deleted")

	return nil
}

// Helper methods

// convertStoreUserToModelsUser converts store.User to models.User
func (s *AuthService) convertStoreUserToModelsUser(storeUser *store.User) models.User {
	return models.User{
		ID:            storeUser.ID,
		Username:      storeUser.Username,
		Email:         storeUser.Email,
		DisplayName:   storeUser.DisplayName,
		AvatarURL:     storeUser.AvatarURL,
		Role:          "viewer", // Will be set by caller based on roles
		IsActive:      storeUser.IsActive,
		EmailVerified: storeUser.EmailVerified,
		LastLogin:     storeUser.LastLoginAt,
		CreatedAt:     storeUser.CreatedAt,
		UpdatedAt:     storeUser.UpdatedAt,
	}
}

// GetUserSessions returns active sessions for a user
func (s *AuthService) GetUserSessions(userID uint) ([]*SessionInfo, error) {
	return s.securityService.GetUserSessions(userID), nil
}

// InvalidateUserSession invalidates a specific user session
func (s *AuthService) InvalidateUserSession(userID uint, sessionID string) error {
	// Verify the session belongs to the user
	sessions := s.securityService.GetUserSessions(userID)
	found := false
	for _, session := range sessions {
		if session.SessionID == sessionID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("session not found or does not belong to user")
	}

	return s.securityService.InvalidateSession(sessionID)
}

// GetUserSecurityInfo returns security events and warnings for a user
func (s *AuthService) GetUserSecurityInfo(userID uint) ([]*store.AuditLog, []string, error) {
	// Get recent audit logs for the user
	events, _, err := s.store.GetAuditLogsByUserID(userID, 0, 50)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Get suspicious activity warnings
	warnings, err := s.securityService.DetectSuspiciousActivity(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to detect suspicious activity: %w", err)
	}

	return events, warnings, nil
}

// ValidatePassword validates a password against security policy
func (s *AuthService) ValidatePassword(password string) (*models.ValidatePasswordResponse, error) {
	validationErrors := s.securityService.ValidatePassword(password)

	// Convert security service errors to model errors
	modelErrors := make([]models.PasswordValidationError, len(validationErrors))
	for i, err := range validationErrors {
		modelErrors[i] = models.PasswordValidationError{
			Field:   err.Field,
			Message: err.Message,
		}
	}

	response := &models.ValidatePasswordResponse{
		Valid:  len(validationErrors) == 0,
		Errors: modelErrors,
	}

	return response, nil
}

// createAuditLog creates an audit log entry
func (s *AuthService) createAuditLog(userID *uint, action, resource, resourceID, ipAddress, userAgent, details string) {
	auditLog := &store.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
	}

	// Don't fail the main operation if audit logging fails
	if err := s.store.CreateAuditLog(auditLog); err != nil {
		fmt.Printf("Failed to create audit log: %v\n", err)
	}
}
