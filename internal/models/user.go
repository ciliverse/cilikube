package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User user model
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:user;size:20"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

//// UserRole user role association table
//type UserRole struct {
//	ID     uint   `gorm:"primaryKey" json:"id"`
//	UserID uint   `gorm:"index" json:"user_id"`
//	Role   string `gorm:"size:50" json:"role"`
//}
//
//// TableName specifies table name
//func (UserRole) TableName() string {
//	return "user_roles"
//}

// LoginRequest
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse returns jwt token after successful login
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}

// HashPassword encrypts password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts to response format
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}

// IsAdmin checks if user is administrator
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// BeforeCreate GORM hook: encrypt password before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.HashPassword()
}

type CreateUserRequest struct {
	Username        string `json:"username"`
	ConfirmPassword string `json:"confirmPassword"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	Roles           string `json:"roles"`
}

type CreateUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
	Status   string `json:"status"`
}
