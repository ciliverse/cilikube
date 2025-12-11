package models

import (
	"time"
)

// OAuthProvider represents OAuth provider information for a user
type OAuthProvider struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	Provider       string     `json:"provider" gorm:"not null;size:50"`
	ProviderUserID string     `json:"provider_user_id" gorm:"not null;size:100"`
	AccessToken    string     `json:"-" gorm:"type:text"`
	RefreshToken   string     `json:"-" gorm:"type:text"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName specifies the table name for OAuthProvider model
func (OAuthProvider) TableName() string {
	return "oauth_providers"
}

// OAuthLoginRequest request for OAuth login
type OAuthLoginRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
	State    string `json:"state" binding:"required"`
}

// OAuthTokenResponse response containing OAuth tokens
type OAuthTokenResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	TokenType    string     `json:"token_type"`
}

// OAuthUserInfo user information from OAuth provider
type OAuthUserInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// OAuthProviderResponse response for OAuth provider operations
type OAuthProviderResponse struct {
	ID             uint       `json:"id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// LinkOAuthAccountRequest request for linking OAuth account
type LinkOAuthAccountRequest struct {
	Provider       string     `json:"provider" binding:"required"`
	ProviderUserID string     `json:"provider_user_id" binding:"required"`
	AccessToken    string     `json:"access_token" binding:"required"`
	RefreshToken   string     `json:"refresh_token"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

// UnlinkOAuthAccountRequest request for unlinking OAuth account
type UnlinkOAuthAccountRequest struct {
	Provider string `json:"provider" binding:"required"`
}

// ToResponse converts OAuthProvider to OAuthProviderResponse
func (o *OAuthProvider) ToResponse() OAuthProviderResponse {
	return OAuthProviderResponse{
		ID:             o.ID,
		Provider:       o.Provider,
		ProviderUserID: o.ProviderUserID,
		ConnectedAt:    o.CreatedAt,
		ExpiresAt:      o.ExpiresAt,
	}
}

// ToProviderInfo converts OAuthProvider to OAuthProviderInfo for user profile
func (o *OAuthProvider) ToProviderInfo() OAuthProviderInfo {
	return OAuthProviderInfo{
		Provider:       o.Provider,
		ProviderUserID: o.ProviderUserID,
		ConnectedAt:    o.CreatedAt,
		ExpiresAt:      o.ExpiresAt,
	}
}

// SupportedOAuthProviders list of supported OAuth providers
var SupportedOAuthProviders = []string{
	"github",
	// Future providers can be added here
	// "google",
	// "gitlab",
}
