package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related requests
type OAuthHandler struct {
	oauthService *service.OAuthService
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService *service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

// GetAuthURL generates OAuth authorization URL
func (h *OAuthHandler) GetAuthURL(c *gin.Context) {
	provider := c.Param("provider")
	state := c.Query("state")

	if state == "" {
		state = "default_state" // In production, generate a secure random state
	}

	authURL, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to generate auth URL",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"auth_url": authURL,
			"state":    state,
		},
		"message": "Auth URL generated successfully",
	})
}

// HandleCallback handles OAuth callback
func (h *OAuthHandler) HandleCallback(c *gin.Context) {
	var req models.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Handle OAuth login
	loginResp, err := h.oauthService.LoginWithOAuth(req.Provider, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "OAuth login failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    loginResp,
		"message": "OAuth login successful",
	})
}

// LinkAccount links OAuth account to current user
func (h *OAuthHandler) LinkAccount(c *gin.Context) {
	// Get current user from JWT token
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Authentication required",
		})
		return
	}

	var req models.OAuthLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Link OAuth account
	if err := h.oauthService.LinkAccount(userID, req.Provider, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to link OAuth account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OAuth account linked successfully",
	})
}

// UnlinkAccount unlinks OAuth account from current user
func (h *OAuthHandler) UnlinkAccount(c *gin.Context) {
	// Get current user from JWT token
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Authentication required",
		})
		return
	}

	var req models.OAuthUnlinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Unlink OAuth account
	if err := h.oauthService.UnlinkAccount(userID, req.Provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Failed to unlink OAuth account",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OAuth account unlinked successfully",
	})
}
