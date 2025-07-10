package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: &service.AuthService{},
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户通过用户名和密码登录系统
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "登录信息"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册账号
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "注册信息"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
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
		"message": "注册成功",
		"data":    response,
	})
}

// GetProfile 获取当前用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料信息
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
			"message": "用户信息不存在",
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
		"message": "获取成功",
		"data":    response,
	})
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.UpdateProfileRequest true "用户资料"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户信息不存在",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
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
		"message": "更新成功",
		"data":    response,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _, _, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户信息不存在",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
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
		"message": "密码修改成功",
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出系统（前端清除token即可）
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}

// GetUserList 获取用户列表（管理员）
// @Summary 获取用户列表
// @Description 管理员获取系统中所有用户列表
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
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
			"message": "获取用户列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"users":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateUserStatus 更新用户状态（管理员）
// @Summary 更新用户状态
// @Description 管理员启用或禁用用户账号
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param status body map[string]bool true "状态信息"
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
			"message": "无效的用户ID",
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	err = h.authService.UpdateUserStatus(uint(userID), req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户状态更新成功",
	})
}

// DeleteUser 删除用户（管理员）
// @Summary 删除用户
// @Description 管理员删除用户账号
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
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
			"message": "无效的用户ID",
		})
		return
	}

	// 防止删除自己
	currentUserID, _, _, ok := auth.GetCurrentUser(c)
	if ok && currentUserID == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不能删除自己的账号",
		})
		return
	}

	err = h.authService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户删除成功",
	})
}
