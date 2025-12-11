package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/utils"
)

// RoleManagementHandler handles role management operations for administrators
type RoleManagementHandler struct {
	roleService *service.RoleService
}

// NewRoleManagementHandler creates a new RoleManagementHandler instance
func NewRoleManagementHandler(roleService *service.RoleService) *RoleManagementHandler {
	return &RoleManagementHandler{
		roleService: roleService,
	}
}

// ListRoles gets all roles in the system
func (h *RoleManagementHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleService.ListRoles()
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get roles", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"roles": roles,
		"total": len(roles),
	}, "Roles retrieved successfully")
}

// GetRole gets a specific role by ID
func (h *RoleManagementHandler) GetRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "Role not found", err.Error())
		return
	}

	// Get users assigned to this role
	users, err := h.roleService.GetRoleUsers(uint(roleID))
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get role users", err.Error())
		return
	}

	response := gin.H{
		"id":           role.ID,
		"name":         role.Name,
		"display_name": role.DisplayName,
		"description":  role.Description,
		"is_system":    role.IsSystem,
		"created_at":   role.CreatedAt,
		"updated_at":   role.UpdatedAt,
		"users":        users,
		"user_count":   len(users),
	}

	utils.ApiSuccess(c, response, "Role retrieved successfully")
}

// CreateRole creates a new role (admin function)
func (h *RoleManagementHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to create role", err.Error())
		return
	}

	utils.ApiSuccess(c, role, "Role created successfully")
}

// UpdateRole updates role information (admin function)
func (h *RoleManagementHandler) UpdateRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(uint(roleID), &req)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to update role", err.Error())
		return
	}

	utils.ApiSuccess(c, role, "Role updated successfully")
}

// DeleteRole deletes a role (admin function)
func (h *RoleManagementHandler) DeleteRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.roleService.DeleteRole(uint(roleID))
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to delete role", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"message": "Role deleted successfully",
		"role_id": roleID,
	}, "Role deleted successfully")
}

// AssignRoleToUser assigns a role to a user (admin function)
func (h *RoleManagementHandler) AssignRoleToUser(c *gin.Context) {
	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Get current user ID for audit
	currentUserID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	err := h.roleService.AssignRoleToUser(req.UserID, req.RoleID, currentUserID)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to assign role", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"message": "Role assigned successfully",
		"user_id": req.UserID,
		"role_id": req.RoleID,
	}, "Role assigned successfully")
}

// RemoveRoleFromUser removes a role from a user (admin function)
func (h *RoleManagementHandler) RemoveRoleFromUser(c *gin.Context) {
	var req models.RemoveRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Get current user ID for audit
	currentUserID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	err := h.roleService.RemoveRoleFromUser(req.UserID, req.RoleID, currentUserID)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to remove role", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"message": "Role removed successfully",
		"user_id": req.UserID,
		"role_id": req.RoleID,
	}, "Role removed successfully")
}

// GetUserRoles gets all roles assigned to a specific user
func (h *RoleManagementHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, err := h.roleService.GetUserRoles(uint(userID))
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"user_id": userID,
		"roles":   roles,
		"total":   len(roles),
	}, "User roles retrieved successfully")
}

// GetRoleUsers gets all users assigned to a specific role
func (h *RoleManagementHandler) GetRoleUsers(c *gin.Context) {
	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	users, err := h.roleService.GetRoleUsers(uint(roleID))
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get role users", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"role_id": roleID,
		"users":   users,
		"total":   len(users),
	}, "Role users retrieved successfully")
}

// GetAvailablePermissions gets all available permissions in the system
func (h *RoleManagementHandler) GetAvailablePermissions(c *gin.Context) {
	// Define available permission categories
	categories := []gin.H{
		{
			"name":         "cluster",
			"display_name": "Cluster Management",
			"description":  "Kubernetes cluster related operation permissions",
			"permissions": []gin.H{
				{"name": "read:clusters", "display_name": "View Clusters", "description": "View cluster information and status"},
				{"name": "write:clusters", "display_name": "Manage Clusters", "description": "Create, update, delete clusters"},
				{"name": "read:nodes", "display_name": "View Nodes", "description": "View cluster node information"},
				{"name": "write:nodes", "display_name": "Manage Nodes", "description": "Manage cluster nodes"},
			},
		},
		{
			"name":         "workload",
			"display_name": "Workloads",
			"description":  "Pod, Deployment and other workload permissions",
			"permissions": []gin.H{
				{"name": "read:pods", "display_name": "View Pods", "description": "View Pod information and logs"},
				{"name": "write:pods", "display_name": "Manage Pods", "description": "Create, update, delete Pods"},
				{"name": "read:deployments", "display_name": "View Deployments", "description": "View deployment information"},
				{"name": "write:deployments", "display_name": "Manage Deployments", "description": "Manage deployments"},
				{"name": "exec:pods", "display_name": "Pod Terminal", "description": "Execute commands in Pods"},
			},
		},
		{
			"name":         "network",
			"display_name": "Network Resources",
			"description":  "Service, Ingress and other network resource permissions",
			"permissions": []gin.H{
				{"name": "read:services", "display_name": "View Services", "description": "View service information"},
				{"name": "write:services", "display_name": "Manage Services", "description": "Manage services"},
				{"name": "read:ingress", "display_name": "View Ingress", "description": "View ingress rules"},
				{"name": "write:ingress", "display_name": "Manage Ingress", "description": "Manage ingress rules"},
			},
		},
		{
			"name":         "storage",
			"display_name": "Storage Resources",
			"description":  "PV, PVC, ConfigMap and other storage permissions",
			"permissions": []gin.H{
				{"name": "read:configmaps", "display_name": "View ConfigMaps", "description": "View configuration maps"},
				{"name": "write:configmaps", "display_name": "Manage ConfigMaps", "description": "Manage configuration maps"},
				{"name": "read:secrets", "display_name": "View Secrets", "description": "View secrets"},
				{"name": "write:secrets", "display_name": "Manage Secrets", "description": "Manage secrets"},
				{"name": "read:pv", "display_name": "View PVs", "description": "View persistent volumes"},
				{"name": "write:pv", "display_name": "Manage PVs", "description": "Manage persistent volumes"},
			},
		},
		{
			"name":         "admin",
			"display_name": "System Administration",
			"description":  "User, role and other system administration permissions",
			"permissions": []gin.H{
				{"name": "admin:users", "display_name": "User Management", "description": "Manage system users", "system_required": true},
				{"name": "admin:roles", "display_name": "Role Management", "description": "Manage system roles", "system_required": true},
				{"name": "admin:system", "display_name": "System Settings", "description": "Manage system configuration", "system_required": true},
				{"name": "admin:audit", "display_name": "Audit Logs", "description": "View audit logs"},
			},
		},
	}

	utils.ApiSuccess(c, gin.H{
		"categories": categories,
	}, "Available permissions retrieved successfully")
}

// GetRolePermissions gets permissions assigned to a specific role
func (h *RoleManagementHandler) GetRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	// For now, return mock permissions
	// TODO: Implement actual permission retrieval from database
	permissions := []string{
		"read:clusters",
		"read:pods",
		"read:deployments",
		"read:services",
	}

	utils.ApiSuccess(c, gin.H{
		"role_id":     roleID,
		"permissions": permissions,
	}, "Role permissions retrieved successfully")
}

// UpdateRolePermissions updates permissions for a specific role
func (h *RoleManagementHandler) UpdateRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// TODO: Implement actual permission update in database
	// For now, just return success

	utils.ApiSuccess(c, gin.H{
		"role_id":     roleID,
		"permissions": req.Permissions,
		"message":     "Role permissions updated successfully",
	}, "Role permissions updated successfully")
}
