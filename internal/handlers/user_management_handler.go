package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
)

// UserManagementHandler handles user management operations for administrators
type UserManagementHandler struct {
	authService *service.AuthService
	roleService *service.RoleService
}

// NewUserManagementHandler creates a new UserManagementHandler instance
func NewUserManagementHandler(authService *service.AuthService, roleService *service.RoleService) *UserManagementHandler {
	return &UserManagementHandler{
		authService: authService,
		roleService: roleService,
	}
}

// ListUsers gets paginated list of users with optional search
func (h *UserManagementHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := strings.TrimSpace(c.Query("search"))
	status := c.Query("status") // active, inactive, all

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Get users from auth service
	users, total, err := h.authService.GetUserList(page, pageSize)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user list", err.Error())
		return
	}

	// Apply search filter if provided
	if search != "" {
		filteredUsers := make([]models.UserResponse, 0)
		searchLower := strings.ToLower(search)

		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Username), searchLower) ||
				strings.Contains(strings.ToLower(user.Email), searchLower) ||
				strings.Contains(strings.ToLower(user.DisplayName), searchLower) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
		total = int64(len(users))
	}

	// Apply status filter if provided
	if status != "" && status != "all" {
		filteredUsers := make([]models.UserResponse, 0)
		isActive := status == "active"

		for _, user := range users {
			if user.IsActive == isActive {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
		total = int64(len(users))
	}

	// Enhance users with role information
	enhancedUsers := make([]gin.H, len(users))
	for i, user := range users {
		roles, err := h.roleService.GetUserRoles(user.ID)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
			return
		}

		roleNames := make([]string, len(roles))
		for j, role := range roles {
			roleNames[j] = role.Name
		}

		enhancedUsers[i] = gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"email":          user.Email,
			"display_name":   user.DisplayName,
			"avatar_url":     user.AvatarURL,
			"is_active":      user.IsActive,
			"email_verified": user.EmailVerified,
			"last_login":     user.LastLogin,
			"created_at":     user.CreatedAt,
			"roles":          roleNames,
		}
	}

	response := gin.H{
		"data": enhancedUsers,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"filters": gin.H{
			"search": search,
			"status": status,
		},
	}

	utils.ApiSuccess(c, response, "User list retrieved successfully")
}

// GetUser gets a specific user by ID
func (h *UserManagementHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get user profile
	profile, err := h.authService.GetProfile(uint(userID))
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Get user roles
	roles, err := h.roleService.GetUserRoles(uint(userID))
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
		return
	}

	// Create enhanced response
	response := gin.H{
		"id":              profile.ID,
		"username":        profile.Username,
		"email":           profile.Email,
		"display_name":    profile.DisplayName,
		"avatar_url":      profile.AvatarURL,
		"is_active":       profile.IsActive,
		"email_verified":  profile.EmailVerified,
		"last_login":      profile.LastLogin,
		"created_at":      profile.CreatedAt,
		"updated_at":      profile.UpdatedAt,
		"roles":           roles,
		"oauth_providers": profile.OAuthProviders,
	}

	utils.ApiSuccess(c, response, "User retrieved successfully")
}

// CreateUser creates a new user (admin function)
func (h *UserManagementHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		utils.ApiError(c, http.StatusBadRequest, "Password and confirm password do not match")
		return
	}

	// Create user using auth service register functionality
	registerReq := models.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := h.authService.Register(&registerReq)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	// Update display name if provided
	if req.DisplayName != "" {
		updateReq := models.UpdateProfileRequest{
			Email:       req.Email,
			DisplayName: req.DisplayName,
			AvatarURL:   "",
		}

		_, err = h.authService.UpdateProfile(createdUser.ID, &updateReq)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, "User created but failed to update display name", err.Error())
			return
		}
	}

	response := models.CreateUserResponse{
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Roles:    req.Roles,
		Status:   "created",
	}

	utils.ApiSuccess(c, response, "User created successfully")
}

// UpdateUser updates user information (admin function)
func (h *UserManagementHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	uid := uint(userID)

	// Check if user exists
	_, err = h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Update user profile
	updateProfileReq := models.UpdateProfileRequest{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
	}

	_, err = h.authService.UpdateProfile(uid, &updateProfileReq)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	// Update user status if provided
	if req.IsActive != nil {
		err = h.authService.UpdateUserStatus(uid, *req.IsActive)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, "Failed to update user status", err.Error())
			return
		}
	}

	// Update user roles if provided
	if len(req.Roles) > 0 {
		// Remove all existing roles
		existingRoles, err := h.roleService.GetUserRoles(uid)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, "Failed to get existing roles", err.Error())
			return
		}

		// Get current admin user ID for audit
		currentUserID := uint(1) // Default admin ID, should be extracted from JWT context
		if userID, exists := c.Get("userID"); exists {
			if uid, ok := userID.(uint); ok {
				currentUserID = uid
			}
		}

		for _, role := range existingRoles {
			if err := h.roleService.RemoveRoleFromUser(uid, role.ID, currentUserID); err != nil {
				utils.ApiError(c, http.StatusInternalServerError, "Failed to remove existing role", err.Error())
				return
			}
		}

		// Assign new roles
		for _, roleName := range req.Roles {
			role, err := h.roleService.GetRoleByName(roleName)
			if err != nil {
				utils.ApiError(c, http.StatusBadRequest, "Invalid role: "+roleName, err.Error())
				return
			}

			if err := h.roleService.AssignRoleToUser(uid, role.ID, currentUserID); err != nil {
				utils.ApiError(c, http.StatusInternalServerError, "Failed to assign role", err.Error())
				return
			}
		}
	}

	// Get updated user profile
	updatedProfile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get updated profile", err.Error())
		return
	}

	// Get updated roles
	roles, err := h.roleService.GetUserRoles(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get updated roles", err.Error())
		return
	}

	response := gin.H{
		"id":              updatedProfile.ID,
		"username":        updatedProfile.Username,
		"email":           updatedProfile.Email,
		"display_name":    updatedProfile.DisplayName,
		"avatar_url":      updatedProfile.AvatarURL,
		"is_active":       updatedProfile.IsActive,
		"email_verified":  updatedProfile.EmailVerified,
		"last_login":      updatedProfile.LastLogin,
		"created_at":      updatedProfile.CreatedAt,
		"updated_at":      updatedProfile.UpdatedAt,
		"roles":           roles,
		"oauth_providers": updatedProfile.OAuthProviders,
	}

	utils.ApiSuccess(c, response, "User updated successfully")
}

// UpdateUserStatus updates user active status
func (h *UserManagementHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	err = h.authService.UpdateUserStatus(uint(userID), req.IsActive)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to update user status", err.Error())
		return
	}

	status := "activated"
	if !req.IsActive {
		status = "deactivated"
	}

	utils.ApiSuccess(c, gin.H{
		"message":   "User " + status + " successfully",
		"user_id":   userID,
		"is_active": req.IsActive,
	}, "User status updated successfully")
}

// DeleteUser deletes a user (admin function)
func (h *UserManagementHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	uid := uint(userID)

	// Prevent self-deletion
	currentUserID, exists := c.Get("userID")
	if exists {
		if cuid, ok := currentUserID.(uint); ok && cuid == uid {
			utils.ApiError(c, http.StatusBadRequest, "Cannot delete your own account")
			return
		}
	}

	// Check if user exists
	_, err = h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Delete user
	err = h.authService.DeleteUser(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{
		"message": "User deleted successfully",
		"user_id": userID,
	}, "User deleted successfully")
}
