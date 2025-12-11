package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
)

// ProfileHandler handles user profile related requests
type ProfileHandler struct {
	authService *service.AuthService
	roleService *service.RoleService
}

// NewProfileHandler creates a new ProfileHandler instance
func NewProfileHandler(authService *service.AuthService, roleService *service.RoleService) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
		roleService: roleService,
	}
}

// GetProfile gets current user profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	profile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get profile", err.Error())
		return
	}

	// Get user roles
	roles, err := h.roleService.GetUserRoles(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
		return
	}

	// Add roles to profile response
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Create enhanced profile response
	enhancedProfile := struct {
		*models.UserProfileResponse
		Roles []string `json:"roles"`
	}{
		UserProfileResponse: profile,
		Roles:               roleNames,
	}

	utils.ApiSuccess(c, enhancedProfile, "Profile retrieved successfully")
}

// UpdateProfile updates user profile information
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(uid, &req)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.ApiSuccess(c, user, "Profile updated successfully")
}

// ChangePassword changes user password
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := h.authService.ChangePassword(uid, &req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{"message": "Password changed successfully"}, "Password changed successfully")
}

// UploadAvatar uploads user avatar
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Failed to get uploaded file", err.Error())
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		utils.ApiError(c, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and GIF are allowed")
		return
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if header.Size > maxSize {
		utils.ApiError(c, http.StatusBadRequest, "File size too large. Maximum size is 5MB")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("avatar_%d_%s%s", uid, uuid.New().String(), ext)

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to create upload directory", err.Error())
		return
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to create file", err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to save file", err.Error())
		return
	}

	// Update user avatar URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	updateReq := models.UpdateProfileRequest{
		AvatarURL: avatarURL,
	}

	// Get current user profile to preserve other fields
	profile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get current profile", err.Error())
		return
	}

	updateReq.Email = profile.Email
	updateReq.DisplayName = profile.DisplayName

	_, err = h.authService.UpdateProfile(uid, &updateReq)
	if err != nil {
		// Clean up uploaded file if profile update fails
		os.Remove(filePath)
		utils.ApiError(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	response := models.AvatarUploadResponse{
		AvatarURL: avatarURL,
		Message:   "Avatar uploaded successfully",
	}

	utils.ApiSuccess(c, response, "Avatar uploaded successfully")
}

// DeleteAvatar removes user avatar
func (h *ProfileHandler) DeleteAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get current user profile
	profile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get current profile", err.Error())
		return
	}

	// Remove avatar file if it exists
	if profile.AvatarURL != "" && strings.HasPrefix(profile.AvatarURL, "/uploads/avatars/") {
		filePath := "." + profile.AvatarURL
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}

	// Update user profile to remove avatar URL
	updateReq := models.UpdateProfileRequest{
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
		AvatarURL:   "", // Clear avatar URL
	}

	_, err = h.authService.UpdateProfile(uid, &updateReq)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.ApiSuccess(c, gin.H{"message": "Avatar deleted successfully"}, "Avatar deleted successfully")
}

// UpdateAvatar updates user avatar URL (for color avatars or preset avatars)
func (h *ProfileHandler) UpdateAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	var req models.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Get current user profile
	profile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user profile", err.Error())
		return
	}

	// Update user avatar URL
	updateReq := models.UpdateProfileRequest{
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
		AvatarURL:   req.AvatarURL,
	}

	updatedProfile, err := h.authService.UpdateProfile(uid, &updateReq)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to update avatar", err.Error())
		return
	}

	response := models.AvatarUploadResponse{
		AvatarURL: updatedProfile.AvatarURL,
		Message:   "Avatar updated successfully",
	}

	utils.ApiSuccess(c, response, "Avatar updated successfully")
}

// GetUserRoles gets current user's roles
func (h *ProfileHandler) GetUserRoles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	roles, err := h.roleService.GetUserRoles(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
		return
	}

	utils.ApiSuccess(c, roles, "User roles retrieved successfully")
}

// GetUserPermissions gets current user's effective permissions
func (h *ProfileHandler) GetUserPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get user roles as a simple permission indicator
	roles, err := h.roleService.GetUserRoles(uid)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "Failed to get user roles", err.Error())
		return
	}

	// Convert roles to simple permission list
	permissions := make([]string, len(roles))
	for i, role := range roles {
		permissions[i] = role.Name
	}

	response := gin.H{
		"roles":       roles,
		"permissions": permissions,
	}

	utils.ApiSuccess(c, response, "User permissions retrieved successfully")
}

// GetActivityLog gets user activity log (placeholder for future implementation)
func (h *ProfileHandler) GetActivityLog(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ApiError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	_, ok := userID.(uint)
	if !ok {
		utils.ApiError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// TODO: Implement actual activity log retrieval from audit service
	// For now, return mock data
	activities := []gin.H{
		{
			"id":        1,
			"action":    "login",
			"resource":  "auth",
			"timestamp": time.Now().Add(-1 * time.Hour),
			"details":   "User logged in successfully",
		},
		{
			"id":        2,
			"action":    "profile_update",
			"resource":  "user",
			"timestamp": time.Now().Add(-2 * time.Hour),
			"details":   "User updated profile information",
		},
	}

	response := gin.H{
		"data": activities,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       len(activities),
			"total_pages": 1,
		},
	}

	utils.ApiSuccess(c, response, "Activity log retrieved successfully")
}
