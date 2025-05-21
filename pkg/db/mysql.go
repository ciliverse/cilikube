package db

import (
	"fmt"
	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitializeDB(cfg *configs.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
		return nil, err
	}
	// 自动迁移表
	if err = db.AutoMigrate(
		&models.User{},
	); err != nil {
		log.Fatalf("迁移数据库失败: %v", err)
	}
	return db, nil
}
