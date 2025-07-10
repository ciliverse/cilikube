package service

import (
	"errors"
	"time"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"gorm.io/gorm"
)

type AuthService struct{}

// Login 用户登录 - 简化版本，不使用数据库
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 简单的硬编码验证，生产环境中应该使用数据库
	if req.Username == "admin" && req.Password == "12345678" {
		// 创建一个模拟的用户对象，正确初始化时间字段
		now := time.Now()
		user := models.User{
			ID:        1,
			Username:  "admin",
			Email:     "admin@cilikube.com",
			Role:      "admin",
			IsActive:  true,
			LastLogin: &now,
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		// 生成JWT token - 使用真实的JWT生成
		token, expiresAt, err := auth.GenerateToken(&user)
		if err != nil {
			return nil, errors.New("生成token失败: " + err.Error())
		}
		
		return &models.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			User:      user.ToResponse(),
		}, nil
	}
	
	return nil, errors.New("用户名或密码错误")
}

// Register 用户注册
func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 创建新用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 密码会在BeforeCreate钩子中加密
		Role:     "user",
		IsActive: true,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetProfile 获取用户资料 - 简化版本，返回符合前端期望的格式
func (s *AuthService) GetProfile(userID uint) (map[string]interface{}, error) {
	// 返回模拟的管理员用户资料
	if userID == 1 {
		return map[string]interface{}{
			"username": "admin",
			"roles":    []string{"admin"}, // 前端期望的是数组格式
		}, nil
	}
	
	return nil, errors.New("用户不存在")
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	var user models.User

	err := database.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 检查邮箱是否被其他用户使用
	var count int64
	database.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, userID).Count(&count)
	if count > 0 {
		return nil, errors.New("邮箱已被其他用户使用")
	}

	// 更新用户信息
	user.Email = req.Email
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User

	err := database.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return errors.New("旧密码错误")
	}

	// 更新密码
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	return database.DB.Save(&user).Error
}

// GetUserList 获取用户列表（管理员功能）
func (s *AuthService) GetUserList(page, pageSize int) ([]models.UserResponse, int64, error) {
	var users []models.User
	var total int64

	// 获取总数
	database.DB.Model(&models.User{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, total, nil
}

// UpdateUserStatus 更新用户状态（管理员功能）
func (s *AuthService) UpdateUserStatus(userID uint, isActive bool) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive).Error
}

// DeleteUser 删除用户（管理员功能）
func (s *AuthService) DeleteUser(userID uint) error {
	return database.DB.Delete(&models.User{}, userID).Error
}
