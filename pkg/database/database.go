package database

import (
	"fmt"
	"log"
	"time"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	if !configs.GlobalConfig.Database.Enabled {
		log.Println("数据库未启用，无需初始化。")
		return nil
	}

	var err error

	dsn := configs.GlobalConfig.GetDSN()

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// // 设置日志级别
	// if configs.GlobalConfig.Server.Mode == "release" {
	// 	gormConfig.Logger = logger.Default.LogMode(logger.Error)
	// }

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	// 首先检查数据库是否启用并且 DB 实例已成功创建
	if !configs.GlobalConfig.Database.Enabled || DB == nil {
		log.Println("数据库未启用或未初始化，跳过迁移。")
		return nil // 不启用或未初始化，不算错误，直接返回
	}

	log.Println("开始数据库自动迁移...") // 添加日志
	err := DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
	return nil
}

// CreateDefaultAdmin 创建默认管理员账户
func CreateDefaultAdmin() error {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	// 如果没有用户，创建默认管理员
	if count == 0 {
		admin := &models.User{
			Username: "admin",
			Email:    "admin@cilikube.com",
			Password: "admin123", // 这个密码会在BeforeCreate钩子中被加密
			Role:     "admin",
			IsActive: true,
		}

		if err := DB.Create(admin).Error; err != nil {
			return fmt.Errorf("failed to create default admin: %v", err)
		}

		log.Println("Default admin user created: username=admin, password=admin123")
	}

	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	// 同样检查 DB 是否为 nil
	if DB == nil {
		log.Println("数据库未初始化，无需关闭。")
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
