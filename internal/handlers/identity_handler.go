package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

type IdentityHandler struct {
	auth *service.AuthService
}

func NewIdentityHandler(authService *service.AuthService) *IdentityHandler {
	return &IdentityHandler{auth: authService}
}

func (h *IdentityHandler) VerifyMFA(c *gin.Context) {
	var req struct {
		MFAToken string `json:"mfa_token" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	resp, err := h.auth.VerifyMFA(req.MFAToken, req.Code, auth.AuditClientIP(c), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "login successful", "data": resp})
}

func (h *IdentityHandler) BeginTOTP(c *gin.Context) {
	userID, _ := c.Get("user_id")
	secret, url, err := h.auth.BeginTOTPSetup(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"secret": secret, "url": url}})
}

func (h *IdentityHandler) ConfirmTOTP(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.auth.ConfirmTOTP(userID.(uint), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "mfa enabled"})
}

func (h *IdentityHandler) DisableTOTP(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.auth.DisableTOTP(userID.(uint), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "mfa disabled"})
}

func (h *IdentityHandler) BeginPasskeyRegister(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID, options, err := h.auth.BeginPasskeyRegister(userID.(uint), c.GetHeader("Origin"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session_id": sessionID, "options": options}})
}

func (h *IdentityHandler) FinishPasskeyRegister(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.auth.FinishPasskeyRegister(userID.(uint), c.Query("session_id"), c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "passkey registered"})
}

func (h *IdentityHandler) BeginPasskeyLogin(c *gin.Context) {
	sessionID, options, err := h.auth.BeginPasskeyLogin(c.GetHeader("Origin"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"session_id": sessionID, "options": options}})
}

func (h *IdentityHandler) FinishPasskeyLogin(c *gin.Context) {
	resp, err := h.auth.FinishPasskeyLogin(c.Query("session_id"), c.Request, auth.AuditClientIP(c), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "login successful", "data": resp})
}

func (h *IdentityHandler) DeletePasskey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	rows, err := h.auth.ListPasskeys(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	for _, row := range rows {
		if uint64(row.ID) == id {
			if err := h.auth.DeletePasskey(userID.(uint), row.CredentialID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "message": "passkey removed"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "passkey not found"})
}

func (h *IdentityHandler) TestLDAP(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	_ = c.ShouldBindJSON(&req)
	msg, err := h.auth.TestLDAP(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"message": msg}})
}

func (h *IdentityHandler) ListPasskeys(c *gin.Context) {
	userID, _ := c.Get("user_id")
	rows, err := h.auth.ListPasskeys(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		out = append(out, gin.H{"id": row.ID, "created_at": row.CreatedAt})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": out})
}
