package models

import (
	"time"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:50"`
	DisplayName string    `json:"display_name" gorm:"not null;size:100"`
	Description string    `json:"description" gorm:"type:text"`
	IsSystem    bool      `json:"is_system" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	RoleID     uint      `json:"role_id" gorm:"not null;index"`
	AssignedBy uint      `json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}

// Permission represents a permission in the system
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Resource    string    `json:"resource" gorm:"not null;size:100;index"`
	Action      string    `json:"action" gorm:"not null;size:50;index"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	RoleID       uint `json:"role_id" gorm:"not null;index"`
	PermissionID uint `json:"permission_id" gorm:"not null;index"`
}

// TableName specifies the table name for RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}

// CreateRoleRequest request for creating a new role
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=50"`
	DisplayName string   `json:"display_name" binding:"required,min=2,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest request for updating a role
type UpdateRoleRequest struct {
	DisplayName string   `json:"display_name" binding:"required,min=2,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Permissions []string `json:"permissions"`
}

// RoleResponse response for role operations
type RoleResponse struct {
	ID              uint                 `json:"id"`
	Name            string               `json:"name"`
	DisplayName     string               `json:"display_name"`
	Description     string               `json:"description"`
	Type            string               `json:"type"`
	IsSystem        bool                 `json:"is_system"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Permissions     []PermissionResponse `json:"permissions,omitempty"`
	MainPermissions []string             `json:"main_permissions"`
	UserCount       int                  `json:"user_count,omitempty"`
	PermissionCount int                  `json:"permission_count,omitempty"`
}

// PermissionResponse response for permission operations
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// AssignRoleRequest request for assigning a role to a user
type AssignRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// RemoveRoleRequest request for removing a role from a user
type RemoveRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// UserRoleResponse response for user role operations
type UserRoleResponse struct {
	UserID     uint           `json:"user_id"`
	Username   string         `json:"username"`
	Roles      []RoleResponse `json:"roles"`
	AssignedAt time.Time      `json:"assigned_at"`
	AssignedBy uint           `json:"assigned_by,omitempty"`
}

// ToResponse converts Role to RoleResponse
func (r *Role) ToResponse() RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// ToResponse converts Permission to PermissionResponse
func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
	}
}

// DefaultRoles defines the system default roles
var DefaultRoles = []Role{
	{
		Name:        "admin",
		DisplayName: "Administrator",
		Description: "Full system access with all permissions",
		IsSystem:    true,
	},
	{
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Can read and write most Kubernetes resources",
		IsSystem:    true,
	},
	{
		Name:        "viewer",
		DisplayName: "Viewer",
		Description: "Read-only access to Kubernetes resources",
		IsSystem:    true,
	},
}
