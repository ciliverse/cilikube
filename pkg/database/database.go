package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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
	dbType := configs.GlobalConfig.Database.Type

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database based on type
	switch dbType {
	case "sqlite":
		// Ensure directory exists for SQLite
		if err := ensureSQLiteDir(dsn); err != nil {
			return fmt.Errorf("failed to create SQLite directory: %v", err)
		}
		DB, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	case "mysql", "":
		// Default to MySQL for backward compatibility
		DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to %s database: %v", dbType, err)
	}

	// Configure connection pool (skip for SQLite)
	if dbType != "sqlite" {
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
	}

	log.Printf("Database connected successfully (type: %s)", dbType)
	return nil
}

// ensureSQLiteDir ensures the directory exists for SQLite database file
func ensureSQLiteDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
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
		// Note: Cluster migration is now handled by the store layer
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
	return nil
}

// CreateDefaultAdmin creates default admin account
func CreateDefaultAdmin() error {
	if DB == nil {
		log.Println("database not initialized, skipping default admin creation.")
		return nil
	}

	var count int64
	DB.Model(&models.User{}).Count(&count)

	// If no users exist, create default admin
	if count == 0 {
		admin := &models.User{
			Username: "admin",
			Email:    "admin@cilikube.com",
			Password: "12345678", // This password will be encrypted in BeforeCreate hook
			Role:     "admin",
			IsActive: true,
		}

		if err := DB.Create(admin).Error; err != nil {
			return fmt.Errorf("failed to create default admin: %v", err)
		}

		log.Println("Default admin user created: username=admin, password=12345678")
	}

	return nil
}

// GetDB returns the database instance
func GetDB() (*gorm.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB, nil
}

// IsDataBaseEnabled returns whether the database is enabled
func IsDataBaseEnabled() bool {
	if configs.GlobalConfig == nil {
		return false
	}
	return configs.GlobalConfig.Database.Enabled
}

// GetDatabaseType returns the current database type
func GetDatabaseType() string {
	if configs.GlobalConfig == nil {
		return "unknown"
	}
	return configs.GlobalConfig.Database.Type
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
