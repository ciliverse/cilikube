package repository

import (
	"errors"

	"github.com/ciliverse/cilikube/api/v1/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(username, password, roles, email string) error {
	user := models.User{
		Username: username,
		Password: password,
		Role:     roles,
		Email:    email,
	}
	if err := r.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
