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

// Login user login - simplified version, does not use database
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Simple hardcoded validation, should use database in production environment
	if req.Username == "admin" && req.Password == "12345678" {
		// Create a mock user object, properly initialize time fields
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

		// Generate JWT token - use real JWT generation
		token, expiresAt, err := auth.GenerateToken(&user)
		if err != nil {
			return nil, errors.New("failed to generate token: " + err.Error())
		}

		return &models.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			User:      user.ToResponse(),
		}, nil
	}

	return nil, errors.New("invalid username or password")
}

// Register user registration
func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// Check if username already exists
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // Password will be encrypted in BeforeCreate hook
		Role:     "user",
		IsActive: true,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetProfile gets user profile - simplified version, returns format expected by frontend
func (s *AuthService) GetProfile(userID uint) (map[string]interface{}, error) {
	// Return mock admin user profile
	if userID == 1 {
		return map[string]interface{}{
			"username": "admin",
			"roles":    []string{"admin"}, // Frontend expects array format
		}, nil
	}

	return nil, errors.New("user does not exist")
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(userID uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	var user models.User

	err := database.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}

	// Check if email is being used by other users
	var count int64
	database.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, userID).Count(&count)
	if count > 0 {
		return nil, errors.New("email is already being used by another user")
	}

	// Update user information
	user.Email = req.Email
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword changes password
func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User

	err := database.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user does not exist")
		}
		return err
	}

	// Verify old password
	if !user.CheckPassword(req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// Update password
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	return database.DB.Save(&user).Error
}

// GetUserList gets user list (admin function)
func (s *AuthService) GetUserList(page, pageSize int) ([]models.UserResponse, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	database.DB.Model(&models.User{}).Count(&total)

	// Paginated query
	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, total, nil
}

// UpdateUserStatus updates user status (admin function)
func (s *AuthService) UpdateUserStatus(userID uint, isActive bool) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive).Error
}

// DeleteUser deletes user (admin function)
func (s *AuthService) DeleteUser(userID uint) error {
	return database.DB.Delete(&models.User{}, userID).Error
}
