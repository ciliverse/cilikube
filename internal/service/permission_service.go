package service

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/internal/store"
)

// PermissionService provides permission management functionality
type PermissionService struct {
	store    store.Store
	enforcer *casbin.Enforcer
}

// NewPermissionService creates a new PermissionService instance
func NewPermissionService(store store.Store, enforcer *casbin.Enforcer) *PermissionService {
	return &PermissionService{
		store:    store,
		enforcer: enforcer,
	}
}

// InitializeDefaultPolicies initializes default RBAC policies for the system
func (s *PermissionService) InitializeDefaultPolicies() error {
	if s.enforcer == nil {
		log.Println("Casbin enforcer not available, skipping policy initialization")
		return nil
	}

	// Define default policies for each role
	defaultPolicies := []struct {
		role   string
		object string
		action string
	}{
		// Admin role - full access
		{"admin", "/api/v1/*", "*"},
		{"admin", "/api/v1/auth/*", "*"},
		{"admin", "/api/v1/roles/*", "*"},
		{"admin", "/api/v1/users/*", "*"},
		{"admin", "/api/v1/clusters/*", "*"},

		// Editor role - read/write access to most resources, but not user/role management
		{"editor", "/api/v1/namespaces/*", "*"},
		{"editor", "/api/v1/pods/*", "*"},
		{"editor", "/api/v1/deployments/*", "*"},
		{"editor", "/api/v1/services/*", "*"},
		{"editor", "/api/v1/configmaps/*", "*"},
		{"editor", "/api/v1/secrets/*", "*"},
		{"editor", "/api/v1/persistentvolumes/*", "*"},
		{"editor", "/api/v1/persistentvolumeclaims/*", "*"},
		{"editor", "/api/v1/ingresses/*", "*"},
		{"editor", "/api/v1/nodes/*", "GET"},
		{"editor", "/api/v1/events/*", "GET"},
		{"editor", "/api/v1/summary/*", "GET"},
		{"editor", "/api/v1/auth/profile", "GET"},
		{"editor", "/api/v1/auth/profile", "PUT"},
		{"editor", "/api/v1/auth/password", "PUT"},

		// Viewer role - read-only access
		{"viewer", "/api/v1/namespaces/*", "GET"},
		{"viewer", "/api/v1/pods/*", "GET"},
		{"viewer", "/api/v1/deployments/*", "GET"},
		{"viewer", "/api/v1/services/*", "GET"},
		{"viewer", "/api/v1/configmaps/*", "GET"},
		{"viewer", "/api/v1/secrets/*", "GET"},
		{"viewer", "/api/v1/persistentvolumes/*", "GET"},
		{"viewer", "/api/v1/persistentvolumeclaims/*", "GET"},
		{"viewer", "/api/v1/ingresses/*", "GET"},
		{"viewer", "/api/v1/nodes/*", "GET"},
		{"viewer", "/api/v1/events/*", "GET"},
		{"viewer", "/api/v1/summary/*", "GET"},
		{"viewer", "/api/v1/auth/profile", "GET"},
		{"viewer", "/api/v1/auth/profile", "PUT"},
		{"viewer", "/api/v1/auth/password", "PUT"},
	}

	// Add policies if they don't exist
	for _, policy := range defaultPolicies {
		if err := s.addPolicyIfNotExists(policy.role, policy.object, policy.action); err != nil {
			return fmt.Errorf("failed to add policy (%s, %s, %s): %w", policy.role, policy.object, policy.action, err)
		}
	}

	// Initialize role inheritance (grouping policies)
	if err := s.initializeRoleInheritance(); err != nil {
		return fmt.Errorf("failed to initialize role inheritance: %w", err)
	}

	log.Println("Default RBAC policies initialized successfully")
	return nil
}

// addPolicyIfNotExists adds a policy if it doesn't already exist
func (s *PermissionService) addPolicyIfNotExists(sub, obj, act string) error {
	has, err := s.enforcer.HasPolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("error checking if policy exists: %w", err)
	}

	if !has {
		added, err := s.enforcer.AddPolicy(sub, obj, act)
		if err != nil {
			return fmt.Errorf("failed to add policy: %w", err)
		}
		if added {
			log.Printf("Successfully added policy: %s, %s, %s", sub, obj, act)
		}
	} else {
		log.Printf("Policy already exists, skipping: %s, %s, %s", sub, obj, act)
	}

	return nil
}

// initializeRoleInheritance sets up role inheritance relationships
func (s *PermissionService) initializeRoleInheritance() error {
	// Get all users and their roles from the store
	users, _, err := s.store.ListUsers(0, 1000) // Get first 1000 users
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	for _, user := range users {
		// Get user roles
		roles, err := s.store.GetUserRoles(user.ID)
		if err != nil {
			log.Printf("Failed to get roles for user %d: %v", user.ID, err)
			continue
		}

		// Add grouping policies for each role
		for _, role := range roles {
			userSubject := fmt.Sprintf("user:%d", user.ID)
			if err := s.addGroupingPolicyIfNotExists(userSubject, role.Name); err != nil {
				log.Printf("Failed to add grouping policy for user %d, role %s: %v", user.ID, role.Name, err)
			}
		}
	}

	return nil
}

// addGroupingPolicyIfNotExists adds a grouping policy if it doesn't already exist
func (s *PermissionService) addGroupingPolicyIfNotExists(user, role string) error {
	has, err := s.enforcer.HasGroupingPolicy(user, role)
	if err != nil {
		return fmt.Errorf("error checking if grouping policy exists: %w", err)
	}

	if !has {
		added, err := s.enforcer.AddGroupingPolicy(user, role)
		if err != nil {
			return fmt.Errorf("failed to add grouping policy: %w", err)
		}
		if added {
			log.Printf("Successfully added grouping policy: %s -> %s", user, role)
		}
	}

	return nil
}

// SyncUserRoles synchronizes user roles with Casbin grouping policies
func (s *PermissionService) SyncUserRoles(userID uint) error {
	if s.enforcer == nil {
		return nil // Skip if Casbin is not available
	}

	userSubject := fmt.Sprintf("user:%d", userID)

	// Remove all existing grouping policies for this user
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, userSubject)
	if err != nil {
		return fmt.Errorf("failed to remove existing grouping policies: %w", err)
	}

	// Get current user roles from store
	roles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return fmt.Errorf("failed to get user roles: %w", err)
	}

	// Add grouping policies for current roles
	for _, role := range roles {
		if err := s.addGroupingPolicyIfNotExists(userSubject, role.Name); err != nil {
			return fmt.Errorf("failed to add grouping policy for role %s: %w", role.Name, err)
		}
	}

	return nil
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (s *PermissionService) CheckPermission(userID uint, object, action string) (bool, error) {
	if s.enforcer == nil {
		// If Casbin is not available, allow all operations (fallback mode)
		log.Printf("Casbin enforcer not available, allowing operation: user %d, %s %s", userID, action, object)
		return true, nil
	}

	userSubject := fmt.Sprintf("user:%d", userID)
	allowed, err := s.enforcer.Enforce(userSubject, object, action)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	log.Printf("Permission check: user %d, %s %s -> %v", userID, action, object, allowed)
	return allowed, nil
}

// AddRolePolicy adds a new policy for a role
func (s *PermissionService) AddRolePolicy(role, object, action string) error {
	if s.enforcer == nil {
		return fmt.Errorf("Casbin enforcer not available")
	}

	return s.addPolicyIfNotExists(role, object, action)
}

// RemoveRolePolicy removes a policy for a role
func (s *PermissionService) RemoveRolePolicy(role, object, action string) error {
	if s.enforcer == nil {
		return fmt.Errorf("Casbin enforcer not available")
	}

	removed, err := s.enforcer.RemovePolicy(role, object, action)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if removed {
		log.Printf("Successfully removed policy: %s, %s, %s", role, object, action)
	} else {
		log.Printf("Policy not found, nothing to remove: %s, %s, %s", role, object, action)
	}

	return nil
}

// GetRolePolicies gets all policies for a role
func (s *PermissionService) GetRolePolicies(role string) ([][]string, error) {
	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer not available")
	}

	policies, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered policy: %w", err)
	}

	return policies, nil
}

// GetUserPermissions gets all effective permissions for a user
func (s *PermissionService) GetUserPermissions(userID uint) ([][]string, error) {
	if s.enforcer == nil {
		return nil, fmt.Errorf("Casbin enforcer not available")
	}

	// Get user roles
	roles, err := s.store.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	var allPermissions [][]string

	// Get permissions for each role
	for _, role := range roles {
		rolePolicies, err := s.GetRolePolicies(role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get policies for role %s: %w", role.Name, err)
		}
		allPermissions = append(allPermissions, rolePolicies...)
	}

	log.Printf("User %d has %d effective permissions", userID, len(allPermissions))
	return allPermissions, nil
}
