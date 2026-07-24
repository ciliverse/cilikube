package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// OAuthService provides OAuth authentication functionality
type OAuthService struct {
	store           store.Store
	config          *configs.Config
	securityService *SecurityService
}

// NewOAuthService creates a new OAuthService instance
func NewOAuthService(store store.Store, config *configs.Config, securityService *SecurityService) *OAuthService {
	if securityService == nil {
		securityService = NewSecurityService(store, config)
	}
	return &OAuthService{
		store:           store,
		config:          config,
		securityService: securityService,
	}
}

// GitHubUserInfo represents user information from GitHub API
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubTokenResponse represents GitHub OAuth token response
type GitHubTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// GetAuthURL generates OAuth authorization URL for the specified provider
func (s *OAuthService) GetAuthURL(provider, state string) (string, error) {
	switch provider {
	case "github":
		if !s.config.OAuth.GitHub.Enabled {
			return "", fmt.Errorf("GitHub OAuth is disabled")
		}
		if strings.TrimSpace(s.config.OAuth.GitHub.ClientID) == "" {
			return "", fmt.Errorf("GitHub OAuth is not configured (missing client_id)")
		}
		return s.getGitHubAuthURL(state), nil
	default:
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}

// ListPublicProviders returns OAuth providers visible on the login page.
func (s *OAuthService) ListPublicProviders() map[string]interface{} {
	ghConfigured := strings.TrimSpace(s.config.OAuth.GitHub.ClientID) != ""
	ghReady := s.config.OAuth.GitHub.Enabled && ghConfigured
	var authURL string
	if ghReady {
		authURL, _ = s.GetAuthURL("github", "cilikube_login")
	}
	return map[string]interface{}{
		"providers": []map[string]interface{}{
			{
				"name":         "github",
				"display_name": "GitHub",
				"enabled":      s.config.OAuth.GitHub.Enabled,
				"configured":   ghConfigured,
				"login_ready":  ghReady,
				"auth_url":     authURL,
			},
		},
		"allow_registration":  s.config.OAuth.AllowRegistration && !k8s.IsShowcase(),
		"auto_link_accounts":  s.config.OAuth.AutoLinkAccounts && !k8s.IsShowcase(),
		"default_oauth_role":  "viewer",
	}
}

// ExchangeToken exchanges authorization code for access token
func (s *OAuthService) ExchangeToken(provider, code string) (*OAuthTokenResponse, error) {
	switch provider {
	case "github":
		return s.exchangeGitHubToken(code)
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}

// GetUserInfo gets user information from OAuth provider
func (s *OAuthService) GetUserInfo(provider, token string) (*OAuthUserInfo, error) {
	switch provider {
	case "github":
		return s.getGitHubUserInfo(token)
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}

// LoginWithOAuth handles OAuth login flow
func (s *OAuthService) LoginWithOAuth(provider, code, ipAddress, userAgent string) (*models.LoginResponse, error) {
	if k8s.IsShowcase() {
		return nil, errors.New("OAuth login is disabled in showcase mode")
	}
	// Exchange code for token
	tokenResp, err := s.ExchangeToken(provider, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.GetUserInfo(provider, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if OAuth provider already exists
	oauthProvider, err := s.store.GetOAuthProviderByProviderUserID(provider, userInfo.ProviderUserID)
	if err == nil {
		// Existing OAuth account, login the associated user
		return s.loginExistingOAuthUser(oauthProvider, tokenResp, ipAddress, userAgent)
	}

	// New OAuth account: optionally auto-link by email
	if s.config.OAuth.AutoLinkAccounts && userInfo.Email != "" {
		existingUser, emailErr := s.store.GetUserByEmail(userInfo.Email)
		if emailErr == nil {
			return s.linkOAuthToExistingUser(existingUser, provider, userInfo, tokenResp, ipAddress, userAgent)
		}
	}

	if !s.config.OAuth.AllowRegistration {
		return nil, errors.New("OAuth registration is disabled; contact an administrator to create an account or enable allow_registration")
	}

	// Create new user with OAuth account
	return s.createNewOAuthUser(provider, userInfo, tokenResp, ipAddress, userAgent)
}

// LinkAccount links an OAuth provider to an existing user account
func (s *OAuthService) LinkAccount(userID uint, provider, code string) error {
	if k8s.IsShowcase() {
		return errors.New("OAuth account linking is disabled in showcase mode")
	}
	// Exchange code for token
	tokenResp, err := s.ExchangeToken(provider, code)
	if err != nil {
		return fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.GetUserInfo(provider, tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if this OAuth account is already linked (store may return non-nil zero value on miss)
	existingProvider, lookupErr := s.store.GetOAuthProviderByProviderUserID(provider, userInfo.ProviderUserID)
	found := lookupErr == nil && existingProvider != nil && existingProvider.ID != 0
	if found && existingProvider.UserID != userID {
		other, uErr := s.store.GetUserByID(existingProvider.UserID)
		if uErr == nil && other != nil {
			return fmt.Errorf("this %s account is already linked to user %q — sign in as that user and unlink, or ask an admin to remove/relink it", provider, other.Username)
		}
		return errors.New("this OAuth account is already linked to another user")
	}

	// Create or update OAuth provider
	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	oauthProvider := &store.OAuthProvider{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderUserID,
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   tokenResp.RefreshToken,
		ExpiresAt:      expiresAt,
	}

	if found {
		oauthProvider.ID = existingProvider.ID
		if err := s.store.UpdateOAuthProvider(oauthProvider); err != nil {
			return fmt.Errorf("failed to update OAuth provider: %w", err)
		}
	} else {
		if err := s.store.CreateOAuthProvider(oauthProvider); err != nil {
			return fmt.Errorf("failed to create OAuth provider: %w", err)
		}
	}

	// Create audit log
	s.createAuditLog(&userID, "oauth_link", "oauth_provider", fmt.Sprintf("%s:%s", provider, userInfo.ProviderUserID), "", "", fmt.Sprintf("OAuth account linked: %s", provider))

	return nil
}

// UnlinkAccount removes OAuth provider from user account
func (s *OAuthService) UnlinkAccount(userID uint, provider string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	linked, err := s.store.ListUserOAuthProviders(userID)
	if err != nil {
		return fmt.Errorf("failed to list OAuth providers: %w", err)
	}

	hasPassword := strings.TrimSpace(user.PasswordHash) != ""
	otherProviders := 0
	for _, p := range linked {
		if p.Provider != provider {
			otherProviders++
		}
	}
	if !hasPassword && otherProviders == 0 {
		return errors.New("cannot unlink the only sign-in method; set a password in Profile first")
	}

	if err := s.store.DeleteOAuthProvider(userID, provider); err != nil {
		return fmt.Errorf("failed to unlink OAuth provider: %w", err)
	}

	// Create audit log
	s.createAuditLog(&userID, "oauth_unlink", "oauth_provider", fmt.Sprintf("%s:%d", provider, userID), "", "", fmt.Sprintf("OAuth account unlinked: %s", provider))

	return nil
}

// GitHub OAuth implementation

func (s *OAuthService) getGitHubAuthURL(state string) string {
	baseURL := "https://github.com/login/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", s.config.OAuth.GitHub.ClientID)
	params.Add("redirect_uri", s.config.OAuth.GitHub.RedirectURL)
	params.Add("scope", "user:email")
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func (s *OAuthService) exchangeGitHubToken(code string) (*OAuthTokenResponse, error) {
	tokenURL := "https://github.com/login/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", s.config.OAuth.GitHub.ClientID)
	data.Set("client_secret", s.config.OAuth.GitHub.ClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub token exchange failed: %s", string(body))
	}

	var tokenResp GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &OAuthTokenResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    0, // GitHub tokens don't expire
	}, nil
}

func (s *OAuthService) getGitHubUserInfo(token string) (*OAuthUserInfo, error) {
	userURL := "https://api.github.com/user"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub user info request failed: %s", string(body))
	}

	var githubUser GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Get user email if not public
	if githubUser.Email == "" {
		email, err := s.getGitHubUserEmail(token)
		if err == nil {
			githubUser.Email = email
		}
	}

	return &OAuthUserInfo{
		ProviderUserID: fmt.Sprintf("%d", githubUser.ID),
		Username:       githubUser.Login,
		Email:          githubUser.Email,
		DisplayName:    githubUser.Name,
		AvatarURL:      githubUser.AvatarURL,
	}, nil
}

func (s *OAuthService) getGitHubUserEmail(token string) (string, error) {
	emailURL := "https://api.github.com/user/emails"

	req, err := http.NewRequest("GET", emailURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub email request failed")
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails: %w", err)
	}

	// Find primary email
	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	// Return first email if no primary found
	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", errors.New("no email found")
}

// Helper methods for OAuth login flow

func (s *OAuthService) loginExistingOAuthUser(oauthProvider *store.OAuthProvider, tokenResp *OAuthTokenResponse, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Get the associated user
	storeUser, err := s.store.GetUserByID(oauthProvider.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !storeUser.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Update OAuth token
	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	oauthProvider.AccessToken = tokenResp.AccessToken
	oauthProvider.RefreshToken = tokenResp.RefreshToken
	oauthProvider.ExpiresAt = expiresAt

	if err := s.store.UpdateOAuthProvider(oauthProvider); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update OAuth token: %v\n", err)
	}

	// Update last login time
	now := time.Now()
	storeUser.LastLoginAt = &now
	if err := s.store.UpdateUser(storeUser); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Convert to models user
	user := s.convertStoreUserToModelsUser(storeUser)

	// Get user roles
	roles, err := s.store.GetUserRoles(storeUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) > 0 {
		user.Role = primaryRoleName(roles)
	} else {
		user.Role = "viewer"
	}

	token, expiresAtJWT, err := s.issueTokenWithSession(&user, storeUser.ID, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Create audit log
	s.createAuditLog(&storeUser.ID, "oauth_login", "user", fmt.Sprintf("%d", storeUser.ID), ipAddress, userAgent, fmt.Sprintf(`{"username":%q,"provider":%q,"ip":%q,"result":"success"}`, storeUser.Username, oauthProvider.Provider, ipAddress))

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAtJWT,
		User:      user.ToResponse(),
	}, nil
}

func (s *OAuthService) linkOAuthToExistingUser(existingUser *store.User, provider string, userInfo *OAuthUserInfo, tokenResp *OAuthTokenResponse, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Check if user is active
	if !existingUser.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Create OAuth provider entry
	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	oauthProvider := &store.OAuthProvider{
		UserID:         existingUser.ID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderUserID,
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   tokenResp.RefreshToken,
		ExpiresAt:      expiresAt,
	}

	if err := s.store.CreateOAuthProvider(oauthProvider); err != nil {
		return nil, fmt.Errorf("failed to create OAuth provider: %w", err)
	}

	// Update user avatar if not set
	if existingUser.AvatarURL == "" && userInfo.AvatarURL != "" {
		existingUser.AvatarURL = userInfo.AvatarURL
		if err := s.store.UpdateUser(existingUser); err != nil {
			// Log error but don't fail login
			fmt.Printf("Failed to update user avatar: %v\n", err)
		}
	}

	// Update last login time
	now := time.Now()
	existingUser.LastLoginAt = &now
	if err := s.store.UpdateUser(existingUser); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Convert to models user
	user := s.convertStoreUserToModelsUser(existingUser)

	// Get user roles
	roles, err := s.store.GetUserRoles(existingUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) > 0 {
		user.Role = primaryRoleName(roles)
	} else {
		user.Role = "viewer"
	}

	token, expiresAtJWT, err := s.issueTokenWithSession(&user, existingUser.ID, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Create audit log
	s.createAuditLog(&existingUser.ID, "oauth_login_link", "user", fmt.Sprintf("%d", existingUser.ID), ipAddress, userAgent, fmt.Sprintf(`{"username":%q,"provider":%q,"ip":%q,"result":"success"}`, existingUser.Username, provider, ipAddress))

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAtJWT,
		User:      user.ToResponse(),
	}, nil
}

func (s *OAuthService) createNewOAuthUser(provider string, userInfo *OAuthUserInfo, tokenResp *OAuthTokenResponse, ipAddress, userAgent string) (*models.LoginResponse, error) {
	if !s.config.OAuth.AllowRegistration {
		return nil, errors.New("OAuth registration is disabled")
	}

	// Create new user
	storeUser := &store.User{
		Username:      userInfo.Username,
		Email:         userInfo.Email,
		DisplayName:   userInfo.DisplayName,
		AvatarURL:     userInfo.AvatarURL,
		IsActive:      true,
		EmailVerified: true, // OAuth emails are considered verified
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

	// Create OAuth provider entry
	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	oauthProvider := &store.OAuthProvider{
		UserID:         storeUser.ID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderUserID,
		AccessToken:    tokenResp.AccessToken,
		RefreshToken:   tokenResp.RefreshToken,
		ExpiresAt:      expiresAt,
	}

	if err := s.store.CreateOAuthProvider(oauthProvider); err != nil {
		return nil, fmt.Errorf("failed to create OAuth provider: %w", err)
	}

	// Convert to models user
	user := s.convertStoreUserToModelsUser(storeUser)
	user.Role = "viewer"

	token, expiresAtJWT, err := s.issueTokenWithSession(&user, storeUser.ID, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Create audit log
	s.createAuditLog(&storeUser.ID, "oauth_register", "user", fmt.Sprintf("%d", storeUser.ID), ipAddress, userAgent, fmt.Sprintf(`{"username":%q,"email":%q,"provider":%q,"ip":%q,"result":"success"}`, storeUser.Username, storeUser.Email, provider, ipAddress))

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAtJWT,
		User:      user.ToResponse(),
		IsNewUser: true,
	}, nil
}

func (s *OAuthService) issueTokenWithSession(user *models.User, userID uint, ipAddress, userAgent string) (string, time.Time, error) {
	if userAgent == "" {
		userAgent = "oauth"
	}
	sessionID, err := s.securityService.CreateSession(userID, ipAddress, userAgent)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create session: %w", err)
	}
	token, expiresAt, err := auth.GenerateToken(user, sessionID)
	if err != nil {
		_ = s.securityService.InvalidateSession(sessionID)
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}
	return token, expiresAt, nil
}

// Helper methods

func (s *OAuthService) convertStoreUserToModelsUser(storeUser *store.User) models.User {
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

func (s *OAuthService) createAuditLog(userID *uint, action, resource, resourceID, ipAddress, userAgent, details string) {
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

// Response types for OAuth operations

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type OAuthUserInfo struct {
	ProviderUserID string `json:"provider_user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	AvatarURL      string `json:"avatar_url"`
}
