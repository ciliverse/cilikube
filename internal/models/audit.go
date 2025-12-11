package models

import (
	"encoding/json"
	"time"
)

// AuditLog represents audit log entries for security and compliance
type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     *uint     `json:"user_id" gorm:"index"`
	Action     string    `json:"action" gorm:"not null;size:100;index"`
	Resource   string    `json:"resource" gorm:"size:100;index"`
	ResourceID string    `json:"resource_id" gorm:"size:100"`
	IPAddress  string    `json:"ip_address" gorm:"size:45"`
	UserAgent  string    `json:"user_agent" gorm:"type:text"`
	Details    string    `json:"details" gorm:"type:json"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogResponse response for audit log operations
type AuditLogResponse struct {
	ID         uint                   `json:"id"`
	UserID     *uint                  `json:"user_id"`
	Username   string                 `json:"username,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Details    map[string]interface{} `json:"details"`
	CreatedAt  time.Time              `json:"created_at"`
}

// CreateAuditLogRequest request for creating audit log entry
type CreateAuditLogRequest struct {
	UserID     *uint                  `json:"user_id"`
	Action     string                 `json:"action" binding:"required"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Details    map[string]interface{} `json:"details"`
}

// AuditLogFilter filter for querying audit logs
type AuditLogFilter struct {
	UserID    *uint      `json:"user_id"`
	Action    string     `json:"action"`
	Resource  string     `json:"resource"`
	IPAddress string     `json:"ip_address"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}

// ToResponse converts AuditLog to AuditLogResponse
func (a *AuditLog) ToResponse() AuditLogResponse {
	var details map[string]interface{}
	if a.Details != "" {
		json.Unmarshal([]byte(a.Details), &details)
	}

	return AuditLogResponse{
		ID:         a.ID,
		UserID:     a.UserID,
		Action:     a.Action,
		Resource:   a.Resource,
		ResourceID: a.ResourceID,
		IPAddress:  a.IPAddress,
		UserAgent:  a.UserAgent,
		Details:    details,
		CreatedAt:  a.CreatedAt,
	}
}

// SetDetails sets the details field from a map
func (a *AuditLog) SetDetails(details map[string]interface{}) error {
	if details == nil {
		a.Details = ""
		return nil
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return err
	}
	a.Details = string(detailsJSON)
	return nil
}

// GetDetails gets the details field as a map
func (a *AuditLog) GetDetails() (map[string]interface{}, error) {
	if a.Details == "" {
		return nil, nil
	}

	var details map[string]interface{}
	err := json.Unmarshal([]byte(a.Details), &details)
	return details, err
}

// Common audit actions
const (
	AuditActionLogin            = "login"
	AuditActionLoginFailed      = "login_failed"
	AuditActionLogout           = "logout"
	AuditActionPasswordChange   = "password_change"
	AuditActionProfileUpdate    = "profile_update"
	AuditActionUserCreate       = "user_create"
	AuditActionUserUpdate       = "user_update"
	AuditActionUserDelete       = "user_delete"
	AuditActionRoleAssign       = "role_assign"
	AuditActionRoleRemove       = "role_remove"
	AuditActionOAuthLink        = "oauth_link"
	AuditActionOAuthUnlink      = "oauth_unlink"
	AuditActionResourceCreate   = "resource_create"
	AuditActionResourceUpdate   = "resource_update"
	AuditActionResourceDelete   = "resource_delete"
	AuditActionPermissionDenied = "permission_denied"
)
