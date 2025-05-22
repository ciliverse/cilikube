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

// Login 用户登录
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User

	// 根据用户名查找用户
	err := database.DB.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	database.DB.Save(&user)

	// 生成JWT token
	token, expiresAt, err := auth.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
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

// GetProfile 获取用户资料
func (s *AuthService) GetProfile(userID uint) (*models.UserResponse, error) {
	var user models.User

	err := database.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
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
