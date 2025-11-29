package database

import (
	"fmt"
	"log"
	"time"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes database connection
func InitDatabase() error {
	if !configs.GlobalConfig.Database.Enabled {
		log.Println("database not enabled, no initialization needed.")
		return nil
	}

	var err error

	dsn := configs.GlobalConfig.GetDSN()

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// // Set log level
	// if configs.GlobalConfig.Server.Mode == "release" {
	// 	gormConfig.Logger = logger.Default.LogMode(logger.Error)
	// }

	// Connect to database
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)           // Set maximum number of connections in idle connection pool
	sqlDB.SetMaxOpenConns(100)          // Set maximum number of open database connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set maximum time a connection can be reused

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate automatically migrates database tables
func AutoMigrate() error {
	// First check if database is enabled and DB instance is successfully created
	if !configs.GlobalConfig.Database.Enabled || DB == nil {
		log.Println("database not enabled or not initialized, skipping migration.")
		return nil // Not enabled or not initialized, not an error, return directly
	}

	log.Println("starting database auto migration...") // Add log
	err := DB.AutoMigrate(
		&models.User{},
		&store.Cluster{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
	return nil
}

// CreateDefaultAdmin creates default admin account
func CreateDefaultAdmin() error {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	// If no users exist, create default admin
	if count == 0 {
		admin := &models.User{
			Username: "admin",
			Email:    "admin@cilikube.com",
			Password: "admin123", // This password will be encrypted in BeforeCreate hook
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

// CloseDatabase closes database connection
func CloseDatabase() error {
	// Also check if DB is nil
	if DB == nil {
		log.Println("database not initialized, no need to close.")
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
