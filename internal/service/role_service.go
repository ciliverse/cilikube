package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
)

// RoleService provides role management functionality
type RoleService struct {
	store             store.Store
	permissionService *PermissionService
}

// NewRoleService creates a new RoleService instance
func NewRoleService(store store.Store) *RoleService {
	return &RoleService{
		store: store,
	}
}

// SetPermissionService sets the permission service for role synchronization
func (s *RoleService) SetPermissionService(permissionService *PermissionService) {
	s.permissionService = permissionService
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.RoleResponse, error) {
	// Check if role name already exists
	_, err := s.store.GetRoleByName(req.Name)
	if err == nil {
		return nil, errors.New("role with this name already exists")
	}

	// Create new role
	role := &store.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsSystem:    false, // User-created roles are not system roles
	}

	if err := s.store.CreateRole(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Create audit log
	s.createAuditLog(nil, "role_create", "role", fmt.Sprintf("%d", role.ID), "", "", fmt.Sprintf("Role '%s' created", role.Name))

	// Convert to response
	response := s.convertStoreRoleToResponse(role)
	return &response, nil
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(roleID uint, req *models.UpdateRoleRequest) (*models.RoleResponse, error) {
	// Get existing role
	role, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Check if it's a system role (system roles cannot be modified)
	if role.IsSystem {
		return nil, errors.New("system roles cannot be modified")
	}

	// Update role fields
	role.DisplayName = req.DisplayName
	role.Description = req.Description

	if err := s.store.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Create audit log
	s.createAuditLog(nil, "role_update", "role", fmt.Sprintf("%d", role.ID), "", "", fmt.Sprintf("Role '%s' updated", role.Name))

	// Convert to response
	response := s.convertStoreRoleToResponse(role)
	return &response, nil
}

// DeleteRole deletes a role
func (s *RoleService) DeleteRole(roleID uint) error {
	// Get existing role
	role, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Check if it's a system role (system roles cannot be deleted)
	if role.IsSystem {
		return errors.New("system roles cannot be deleted")
	}

	// Check if role is assigned to any users
	users, err := s.store.GetRoleUsers(roleID)
	if err != nil {
		return fmt.Errorf("failed to check role assignments: %w", err)
	}

	if len(users) > 0 {
		return fmt.Errorf("cannot delete role: it is assigned to %d user(s)", len(users))
	}

	// Delete role
	if err := s.store.DeleteRole(roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Create audit log
	s.createAuditLog(nil, "role_delete", "role", fmt.Sprintf("%d", roleID), "", "", fmt.Sprintf("Role '%s' deleted", role.Name))

	return nil
}

// GetRole gets a role by ID
func (s *RoleService) GetRole(roleID uint) (*models.RoleResponse, error) {
	role, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Get user count for this role
	users, err := s.store.GetRoleUsers(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	response := s.convertStoreRoleToResponse(role)
	response.UserCount = len(users)
	return &response, nil
}

// GetRoleByName gets a role by name
func (s *RoleService) GetRoleByName(name string) (*models.RoleResponse, error) {
	role, err := s.store.GetRoleByName(name)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Get user count for this role
	users, err := s.store.GetRoleUsers(role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	response := s.convertStoreRoleToResponse(role)
	response.UserCount = len(users)
	return &response, nil
}

// ListRoles gets all roles
func (s *RoleService) ListRoles() ([]models.RoleResponse, error) {
	roles, err := s.store.ListRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	responses := make([]models.RoleResponse, len(roles))
	for i, role := range roles {
		// Get user count for each role
		users, err := s.store.GetRoleUsers(role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role users: %w", err)
		}

		responses[i] = s.convertStoreRoleToResponse(role)
		responses[i].UserCount = len(users)
	}

	return responses, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID, roleID uint, assignedBy uint) error {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if role exists
	role, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Check if user already has this role
	hasRole, err := s.store.HasRole(userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if hasRole {
		return nil // Already assigned, no error
	}

	// Assign role
	if err := s.store.AssignRole(userID, roleID); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Sync user roles with Casbin if permission service is available
	if s.permissionService != nil {
		if err := s.permissionService.SyncUserRoles(userID); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to sync user roles with Casbin: %v\n", err)
		}
	}

	// Create audit log
	s.createAuditLog(&assignedBy, "role_assign", "user_role", fmt.Sprintf("%d_%d", userID, roleID), "", "",
		fmt.Sprintf("Role '%s' assigned to user %d", role.Name, userID))

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *RoleService) RemoveRoleFromUser(userID, roleID uint, removedBy uint) error {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if role exists
	role, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Check if user has this role
	hasRole, err := s.store.HasRole(userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if !hasRole {
		return errors.New("user does not have this role")
	}

	// Remove role
	if err := s.store.RemoveRole(userID, roleID); err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	// Sync user roles with Casbin if permission service is available
	if s.permissionService != nil {
		if err := s.permissionService.SyncUserRoles(userID); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to sync user roles with Casbin: %v\n", err)
		}
	}

	// Create audit log
	s.createAuditLog(&removedBy, "role_remove", "user_role", fmt.Sprintf("%d_%d", userID, roleID), "", "",
		fmt.Sprintf("Role '%s' removed from user %d", role.Name, userID))

	return nil
}

// AssignRolesToUser assigns multiple roles to a user (replaces existing roles)
func (s *RoleService) AssignRolesToUser(userID uint, roleIDs []uint, assignedBy uint) error {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate all roles exist
	for _, roleID := range roleIDs {
		_, err := s.store.GetRoleByID(roleID)
		if err != nil {
			return fmt.Errorf("role with ID %d not found", roleID)
		}
	}

	// Get current user roles
	currentRoles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return fmt.Errorf("failed to get current user roles: %w", err)
	}

	// Create maps for easier comparison
	currentRoleIDs := make(map[uint]bool)
	for _, role := range currentRoles {
		currentRoleIDs[role.ID] = true
	}

	newRoleIDs := make(map[uint]bool)
	for _, roleID := range roleIDs {
		newRoleIDs[roleID] = true
	}

	// Remove roles that are no longer assigned
	for _, role := range currentRoles {
		if !newRoleIDs[role.ID] {
			if err := s.store.RemoveRole(userID, role.ID); err != nil {
				return fmt.Errorf("failed to remove role %d: %w", role.ID, err)
			}
			// Create audit log for removal
			s.createAuditLog(&assignedBy, "role_remove", "user_role", fmt.Sprintf("%d_%d", userID, role.ID), "", "",
				fmt.Sprintf("Role '%s' removed from user %d", role.Name, userID))
		}
	}

	// Add new roles
	for _, roleID := range roleIDs {
		if !currentRoleIDs[roleID] {
			if err := s.store.AssignRole(userID, roleID); err != nil {
				return fmt.Errorf("failed to assign role %d: %w", roleID, err)
			}
			// Get role name for audit log
			role, _ := s.store.GetRoleByID(roleID)
			roleName := fmt.Sprintf("%d", roleID)
			if role != nil {
				roleName = role.Name
			}
			// Create audit log for assignment
			s.createAuditLog(&assignedBy, "role_assign", "user_role", fmt.Sprintf("%d_%d", userID, roleID), "", "",
				fmt.Sprintf("Role '%s' assigned to user %d", roleName, userID))
		}
	}

	// Sync user roles with Casbin if permission service is available
	if s.permissionService != nil {
		if err := s.permissionService.SyncUserRoles(userID); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to sync user roles with Casbin: %v\n", err)
		}
	}

	return nil
}

// GetUserRoles gets all roles assigned to a user
func (s *RoleService) GetUserRoles(userID uint) ([]models.RoleResponse, error) {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	roles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	responses := make([]models.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = s.convertStoreRoleToResponse(role)
	}

	return responses, nil
}

// GetRoleUsers gets all users assigned to a role
func (s *RoleService) GetRoleUsers(roleID uint) ([]models.UserRoleResponse, error) {
	// Check if role exists
	_, err := s.store.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	users, err := s.store.GetRoleUsers(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	responses := make([]models.UserRoleResponse, len(users))
	for i, user := range users {
		// Get all roles for this user
		userRoles, err := s.store.GetUserRoles(user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user roles: %w", err)
		}

		roleResponses := make([]models.RoleResponse, len(userRoles))
		for j, userRole := range userRoles {
			roleResponses[j] = s.convertStoreRoleToResponse(userRole)
		}

		responses[i] = models.UserRoleResponse{
			UserID:     user.ID,
			Username:   user.Username,
			Roles:      roleResponses,
			AssignedAt: time.Now(), // TODO: Get actual assignment time from UserRole table
		}
	}

	return responses, nil
}

// InitializeDefaultRoles creates the default system roles if they don't exist
func (s *RoleService) InitializeDefaultRoles() error {
	defaultRoles := []store.Role{
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "System administrator with all permissions",
			IsSystem:    true,
		},
		{
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Can edit most resources, but cannot manage users and clusters",
			IsSystem:    true,
		},
		{
			Name:        "viewer",
			DisplayName: "Viewer",
			Description: "Can only view resources, cannot perform modification operations",
			IsSystem:    true,
		},
	}

	for _, defaultRole := range defaultRoles {
		// Check if role already exists
		_, err := s.store.GetRoleByName(defaultRole.Name)
		if err == nil {
			// Role already exists, skip
			continue
		}

		// Create the role
		role := defaultRole // Create a copy
		if err := s.store.CreateRole(&role); err != nil {
			return fmt.Errorf("failed to create default role '%s': %w", defaultRole.Name, err)
		}

		// Create audit log
		s.createAuditLog(nil, "role_create", "role", fmt.Sprintf("%d", role.ID), "", "",
			fmt.Sprintf("Default system role '%s' initialized", role.Name))
	}

	return nil
}

// Helper methods

// convertStoreRoleToResponse converts store.Role to models.RoleResponse
func (s *RoleService) convertStoreRoleToResponse(role *store.Role) models.RoleResponse {
	// Determine role type
	roleType := "custom"
	if role.IsSystem {
		roleType = "system"
	}

	// Mock main permissions for now - in a real implementation,
	// you would fetch these from the database
	mainPermissions := []string{}
	if role.Name == "admin" {
		mainPermissions = []string{"admin:users", "admin:roles", "admin:system"}
	} else if role.Name == "editor" {
		mainPermissions = []string{"read:clusters", "write:pods", "write:deployments"}
	} else if role.Name == "viewer" {
		mainPermissions = []string{"read:clusters", "read:pods", "read:deployments"}
	}

	return models.RoleResponse{
		ID:              role.ID,
		Name:            role.Name,
		DisplayName:     role.DisplayName,
		Description:     role.Description,
		Type:            roleType,
		IsSystem:        role.IsSystem,
		CreatedAt:       role.CreatedAt,
		UpdatedAt:       role.UpdatedAt,
		MainPermissions: mainPermissions,
		PermissionCount: len(mainPermissions),
	}
}

// createAuditLog creates an audit log entry
func (s *RoleService) createAuditLog(userID *uint, action, resource, resourceID, ipAddress, userAgent, details string) {
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
