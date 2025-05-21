package dao

import "gorm.io/gorm"

type AuthDao struct {
	db *gorm.DB
}

func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{
		db: db,
	}
}
