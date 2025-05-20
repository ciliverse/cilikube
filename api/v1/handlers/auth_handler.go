package handlers

import (
	"errors"
	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

var ErrUserNotFound = errors.New("用户不存在")
var ErrInvalidPassword = errors.New("密码错误")
var ErrUserAlreadyExists = errors.New("用户名已存在")

func (h *AuthHandler) Login(c *gin.Context) {

	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		respondError(c, http.StatusBadRequest, "系统错误, 请求参数解析失败（请检查你的json格式)")
		return
	}
	username := strings.TrimSpace(loginReq.Username)
	password := strings.TrimSpace(loginReq.Password)

	u, err := h.authService.Login(&models.User{
		Username: username,
		Password: password,
	})

	if errors.Is(err, ErrUserNotFound) {
		respondError(c, http.StatusBadRequest, "用户不存在")
		return
	}
	if errors.Is(err, ErrInvalidPassword) {
		respondError(c, http.StatusBadRequest, "密码错误")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "系统错误, 登录失败")
		return
	}
	jwtSecret := os.Getenv("CILIKUBE_JWT_SECRET")
	token, err := auth.GenerateJWT(u.Username, u.ID, u.Roles, jwtSecret)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "系统错误, 登录失败")
		return
	}
	c.Header("cilikube-token", token)

	respondSuccess(c, http.StatusOK, models.LoginResponse{
		Username: u.Username,
		Roles:    u.Roles,
		Status:   "登录成功",
		Token:    token,
	})

	return
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var cu models.CreateUserRequest
	if err := c.ShouldBindJSON(&cu); err != nil {
		respondError(c, http.StatusBadRequest, "系统错误, 请求参数解析失败（请检查你的json格式)")
		return
	}

	// 下面要对输入的两次密码进行校验，如果密码不一致，则返回错误
	if cu.Password != cu.ConfirmPassword {
		respondError(c, http.StatusBadRequest, "两次输入的密码不一致，请检查你输出的密码!")
		return
	}

	// --- 本次没有对密码进行校验操作，需要后续使用的时候添加密码格式的校验工作 ---

	user := &models.User{
		Username: cu.Username,
		Password: cu.Password,
		Email:    cu.Email,
		Roles:    cu.Roles,
	}

	if err := h.authService.CreateUser(user); errors.Is(err, ErrUserAlreadyExists) {
		respondError(c, http.StatusBadRequest, "用户名已存在！")
		return
	} else if err != nil {
		respondError(c, http.StatusInternalServerError, "系统错误, 创建用户失败")
		return
	}

	respondSuccess(c, http.StatusOK, models.CreateUserResponse{
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		Status:   "创建用户成功",
	})
	return
}
