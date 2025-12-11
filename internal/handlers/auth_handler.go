package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login user login
// @Summary User login
// @Description User logs into the system with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login information"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	// Get client IP and user agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "login successful",
		"data":    response,
	})
}

// Register user registration
// @Summary User registration
// @Description New user registers an account
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "Registration information"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "registration successful",
		"data":    response,
	})
}

// GetProfile gets current user profile (legacy format for backward compatibility)
// @Summary Get user profile
// @Description Get profile information of currently logged in user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	response, err := h.authService.GetProfileLegacy(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "retrieved successfully",
		"data":    response,
	})
}

// GetDetailedProfile gets detailed current user profile
// @Summary Get detailed user profile
// @Description Get detailed profile information of currently logged in user including OAuth providers
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/profile/detailed [get]
func (h *AuthHandler) GetDetailedProfile(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	response, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "retrieved successfully",
		"data":    response,
	})
}

// RefreshToken refreshes JWT token
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token that is close to expiry
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Authorization header is required",
		})
		return
	}

	// Extract token from "Bearer <token>"
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid authorization header format",
		})
		return
	}

	response, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "token refreshed successfully",
		"data":    response,
	})
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update profile information of currently logged in user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.UpdateProfileRequest true "User profile"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	response, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "updated successfully",
		"data":    response,
	})
}

// ChangePassword changes password
// @Summary Change password
// @Description Change password of currently logged in user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Password information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "password changed successfully",
	})
}

// Logout user logout
// @Summary User logout
// @Description User logs out of the system and invalidates session
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if ok {
		// Invalidate user sessions
		if err := h.authService.Logout(userID); err != nil {
			// Log error but don't fail logout
			// Frontend will clear token regardless
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "logout successful",
	})
}

// GetUserSessions gets current user's active sessions
// @Summary Get user sessions
// @Description Get list of active sessions for current user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/sessions [get]
func (h *AuthHandler) GetUserSessions(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	sessions, err := h.authService.GetUserSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get user sessions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "retrieved successfully",
		"data":    sessions,
	})
}

// InvalidateSession invalidates a specific session
// @Summary Invalidate session
// @Description Invalidate a specific user session
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/sessions/{sessionId} [delete]
func (h *AuthHandler) InvalidateSession(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "session ID is required",
		})
		return
	}

	err := h.authService.InvalidateUserSession(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "session invalidated successfully",
	})
}

// GetSecurityEvents gets security events for current user
// @Summary Get security events
// @Description Get security events and suspicious activity for current user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/security/events [get]
func (h *AuthHandler) GetSecurityEvents(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "user information does not exist",
		})
		return
	}

	events, warnings, err := h.authService.GetUserSecurityInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get security events: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "retrieved successfully",
		"data": gin.H{
			"events":   events,
			"warnings": warnings,
		},
	})
}

// ValidatePassword validates password against security policy
// @Summary Validate password
// @Description Validate password against current security policy
// @Tags Auth
// @Accept json
// @Produce json
// @Param password body models.ValidatePasswordRequest true "Password to validate"
// @Success 200 {object} models.ValidatePasswordResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/validate-password [post]
func (h *AuthHandler) ValidatePassword(c *gin.Context) {
	var req models.ValidatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	response, err := h.authService.ValidatePassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to validate password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "validation completed",
		"data":    response,
	})
}

// GetUserList gets user list (admin)
// @Summary Get user list
// @Description Admin gets list of all users in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/auth/users [get]
func (h *AuthHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.authService.GetUserList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get user list: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "retrieved successfully",
		"data": gin.H{
			"users":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateUserStatus updates user status (admin)
// @Summary Update user status
// @Description Admin enables or disables user account
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param status body map[string]bool true "Status information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/auth/users/{id}/status [put]
func (h *AuthHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid user ID",
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "parameter error: " + err.Error(),
		})
		return
	}

	err = h.authService.UpdateUserStatus(uint(userID), req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to update user status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "user status updated successfully",
	})
}

// DeleteUser deletes user (admin)
// @Summary Delete user
// @Description Admin deletes user account
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/auth/users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid user ID",
		})
		return
	}

	// Prevent deleting oneself
	currentUserID, _, _, ok := auth.GetCurrentUser(c)
	if ok && currentUserID == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "cannot delete your own account",
		})
		return
	}

	err = h.authService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "user deleted successfully",
	})
}
