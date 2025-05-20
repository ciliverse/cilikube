package service

import (
	"errors"
	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

var ErrUserNotFound = errors.New("用户不存在")
var ErrInvalidPassword = errors.New("密码错误")
var ErrUserAlreadyExists = errors.New("用户名已存在")

func (s *AuthService) Login(u *models.User) (*models.User, error) {
	// 先寻找用户
	user, err := s.repo.FindUserByUsername(u.Username)
	// 用户不存在判断
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	// 其他错误
	if err != nil {
		return nil, err
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return nil, ErrInvalidPassword
	}
	return user, nil
}

func (s *AuthService) CreateUser(u *models.User) error {
	// 检查用户名是否已存在
	_, err := s.repo.FindUserByUsername(u.Username)
	if err == nil {
		log.Println("用户名已存在，跳过创建。")
		return ErrUserAlreadyExists
	}

	// 加密用户的密码
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("加密密码失败: %v", err)
		return err
	}
	err = s.repo.CreateUser(u.Username, string(hashPassword), u.Roles, u.Email)
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
		return err
	}

	return nil

}
