package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex" json:"username"`
	Password  string         `gorm:"size:100" json:"-"` // 不返回密码
	Email     string         `gorm:"size:100" json:"email"`
	Roles     string         `gorm:"size:100" json:"roles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

//// UserRole 用户角色关联表
//type UserRole struct {
//	ID     uint   `gorm:"primaryKey" json:"id"`
//	UserID uint   `gorm:"index" json:"user_id"`
//	Role   string `gorm:"size:50" json:"role"`
//}
//
//// TableName 指定表名
//func (UserRole) TableName() string {
//	return "user_roles"
//}

// LoginRequest
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录成功后返回 jwt token
type LoginResponse struct {
	Username string `json:"username"`
	Roles    string `json:"roles"`
	Status   string `json:"status"`
	Token    string `json:"token"`
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
