package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User user model
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password      string         `json:"-" gorm:"not null"`
	DisplayName   string         `json:"display_name" gorm:"size:100"`
	AvatarURL     string         `json:"avatar_url" gorm:"size:10000"`
	Role          string         `json:"role" gorm:"default:user;size:20"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	LastLogin     *time.Time     `json:"last_login"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

//// UserRole user role association table
//type UserRole struct {
//	ID     uint   `gorm:"primaryKey" json:"id"`
//	UserID uint   `gorm:"index" json:"user_id"`
//	Role   string `gorm:"size:50" json:"role"`
//}
//
//// TableName specifies table name
//func (UserRole) TableName() string {
//	return "user_roles"
//}

// LoginRequest
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse returns jwt token after successful login
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"display_name" binding:"max=100"`
	AvatarURL   string `json:"avatar_url" binding:"max=10000"`
}

type UserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	DisplayName   string     `json:"display_name"`
	AvatarURL     string     `json:"avatar_url"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLogin     *time.Time `json:"last_login"`
	CreatedAt     time.Time  `json:"created_at"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OAuth related request/response types
type OAuthLinkRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type OAuthUnlinkRequest struct {
	Provider string `json:"provider" binding:"required"`
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}

// HashPassword encrypts password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts to response format
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     u.AvatarURL,
		Role:          u.Role,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		LastLogin:     u.LastLogin,
		CreatedAt:     u.CreatedAt,
	}
}

// IsAdmin checks if user is administrator
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// BeforeCreate GORM hook: encrypt password before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword()
}

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6"`
	Password        string `json:"password" binding:"required,min=6"`
	Email           string `json:"email" binding:"required,email"`
	DisplayName     string `json:"display_name" binding:"max=100"`
	Roles           string `json:"roles"`
}

type CreateUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
	Status   string `json:"status"`
}

type UpdateUserRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	DisplayName string   `json:"display_name" binding:"max=100"`
	AvatarURL   string   `json:"avatar_url" binding:"max=10000"`
	IsActive    *bool    `json:"is_active,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// UpdateAvatarRequest request for updating avatar URL
type UpdateAvatarRequest struct {
	AvatarURL    string                 `json:"avatar_url" binding:"required"`
	AvatarType   string                 `json:"avatar_type,omitempty"`
	AvatarConfig map[string]interface{} `json:"avatar_config,omitempty"`
}

// AvatarUploadResponse response for avatar upload
type AvatarUploadResponse struct {
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}

// UserProfileResponse detailed user profile response
type UserProfileResponse struct {
	ID             uint                `json:"id"`
	Username       string              `json:"username"`
	Email          string              `json:"email"`
	DisplayName    string              `json:"display_name"`
	AvatarURL      string              `json:"avatar_url"`
	Role           string              `json:"role"`
	Roles          []string            `json:"roles"`
	IsActive       bool                `json:"is_active"`
	EmailVerified  bool                `json:"email_verified"`
	LastLogin      *time.Time          `json:"last_login"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	OAuthProviders []OAuthProviderInfo `json:"oauth_providers,omitempty"`
}

// OAuthProviderInfo basic OAuth provider information for user profile
type OAuthProviderInfo struct {
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// ToProfileResponse converts to detailed profile response format
func (u *User) ToProfileResponse() UserProfileResponse {
	return UserProfileResponse{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     u.AvatarURL,
		Role:          u.Role,
		Roles:         []string{}, // Will be filled by service layer
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		LastLogin:     u.LastLogin,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// Security-related request/response types

// ValidatePasswordRequest request for password validation
type ValidatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// PasswordValidationError represents a password validation error
type PasswordValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidatePasswordResponse response for password validation
type ValidatePasswordResponse struct {
	Valid  bool                      `json:"valid"`
	Errors []PasswordValidationError `json:"errors,omitempty"`
}

// SessionInfo represents session information for API responses
type SessionInfo struct {
	SessionID string    `json:"session_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SecurityEventResponse represents a security event for API responses
type SecurityEventResponse struct {
	ID        uint      `json:"id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// UserSecurityInfoResponse represents user security information
type UserSecurityInfoResponse struct {
	Events   []SecurityEventResponse `json:"events"`
	Warnings []string                `json:"warnings"`
}
