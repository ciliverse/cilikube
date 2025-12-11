package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server" json:"server"`
	Kubernetes KubernetesConfig `yaml:"kubernetes" json:"kubernetes"`
	Installer  InstallerConfig  `yaml:"installer" json:"installer"`
	Database   DatabaseConfig   `yaml:"database" json:"database"`
	Storage    StorageConfig    `yaml:"storage" json:"storage"`
	JWT        JWTConfig        `yaml:"jwt" json:"jwt"`
	OAuth      OAuthConfig      `yaml:"oauth" json:"oauth"`
	Security   SecurityConfig   `yaml:"security" json:"security"`
	Clusters   []ClusterInfo    `yaml:"clusters" json:"clusters"`
}

type ServerConfig struct {
	Port            string `yaml:"port" json:"port"`
	ReadTimeout     int    `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout" json:"write_timeout"`
	Mode            string `yaml:"mode" json:"mode"`                   // debug, release
	ActiveClusterID string `yaml:"activeCluster" json:"activeCluster"` // Modified to match field name in config file
	EncryptionKey   string `yaml:"encryptionKey" json:"encryptionKey"`
}

type KubernetesConfig struct {
	Kubeconfig string `yaml:"kubeconfig" json:"kubeconfig"`
}

type InstallerConfig struct {
	MinikubePath   string `yaml:"minikubePath" json:"minikubePath"`
	MinikubeDriver string `yaml:"minikubeDriver" json:"minikubeDriver"`
	DownloadDir    string `yaml:"downloadDir" json:"downloadDir"`
}

type DatabaseConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Type     string `yaml:"type" json:"type"` // "mysql", "postgresql", "sqlite"
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"` // Ensure this is username
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"` // Ensure this is database
	Charset  string `yaml:"charset" json:"charset"`
}

type StorageConfig struct {
	Type     string          `yaml:"type" json:"type"` // "memory" or "database", optional, automatically determined based on database configuration by default
	Database *DatabaseConfig `yaml:"database" json:"database"`
}

type OAuthConfig struct {
	GitHub GitHubOAuthConfig `yaml:"github" json:"github"`
}

type GitHubOAuthConfig struct {
	ClientID     string `yaml:"client_id" json:"client_id"`
	ClientSecret string `yaml:"client_secret" json:"client_secret"`
	RedirectURL  string `yaml:"redirect_url" json:"redirect_url"`
}

type JWTConfig struct {
	SecretKey      string        `yaml:"secret_key" json:"secret_key"`
	ExpireDuration time.Duration `yaml:"expire_duration" json:"expire_duration"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
}

type SecurityConfig struct {
	Password    PasswordConfig    `yaml:"password" json:"password"`
	AccountLock AccountLockConfig `yaml:"account_lock" json:"account_lock"`
	Session     SessionConfig     `yaml:"session" json:"session"`
	RateLimit   RateLimitConfig   `yaml:"rate_limit" json:"rate_limit"`
}

type PasswordConfig struct {
	MinLength        int  `yaml:"min_length" json:"min_length"`
	RequireUppercase bool `yaml:"require_uppercase" json:"require_uppercase"`
	RequireLowercase bool `yaml:"require_lowercase" json:"require_lowercase"`
	RequireNumbers   bool `yaml:"require_numbers" json:"require_numbers"`
	RequireSymbols   bool `yaml:"require_symbols" json:"require_symbols"`
	MaxAge           int  `yaml:"max_age" json:"max_age"` // days, 0 means no expiration
}

type AccountLockConfig struct {
	Enabled           bool          `yaml:"enabled" json:"enabled"`
	MaxFailedAttempts int           `yaml:"max_failed_attempts" json:"max_failed_attempts"`
	LockoutDuration   time.Duration `yaml:"lockout_duration" json:"lockout_duration"`
	ResetAfter        time.Duration `yaml:"reset_after" json:"reset_after"` // Reset failed attempts counter after this duration
}

type SessionConfig struct {
	MaxConcurrentSessions int           `yaml:"max_concurrent_sessions" json:"max_concurrent_sessions"`
	IdleTimeout           time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	AbsoluteTimeout       time.Duration `yaml:"absolute_timeout" json:"absolute_timeout"`
	RequireReauth         bool          `yaml:"require_reauth" json:"require_reauth"` // Require re-authentication for sensitive operations
}

type RateLimitConfig struct {
	Enabled       bool          `yaml:"enabled" json:"enabled"`
	LoginAttempts int           `yaml:"login_attempts" json:"login_attempts"` // Max login attempts per window
	LoginWindow   time.Duration `yaml:"login_window" json:"login_window"`     // Time window for login attempts
	APIRequests   int           `yaml:"api_requests" json:"api_requests"`     // Max API requests per window
	APIWindow     time.Duration `yaml:"api_window" json:"api_window"`         // Time window for API requests
	BurstSize     int           `yaml:"burst_size" json:"burst_size"`         // Allow burst requests
}

type ClusterInfo struct {
	// ID is the unique identifier for the cluster, using UUID format
	// If empty, the system will automatically generate a UUID
	ID string `yaml:"id" json:"id"`

	// Name is the user-friendly display name for the cluster
	Name string `yaml:"name" json:"name"`

	// ConfigPath can be the absolute path to kubeconfig file, or "in-cluster"
	ConfigPath string `yaml:"config_path" json:"config_path"`

	// Description cluster description information
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// Provider cloud service provider or environment type, such as "aws", "gcp", "minikube", "on-premise"
	Provider string `yaml:"provider,omitempty" json:"provider,omitempty"`

	// Environment environment identifier, such as "production", "staging", "development"
	Environment string `yaml:"environment,omitempty" json:"environment,omitempty"`

	// Region the region where the cluster is located
	Region string `yaml:"region,omitempty" json:"region,omitempty"`

	// IsActive marks whether this cluster configuration is enabled
	IsActive bool `yaml:"is_active" json:"is_active"`

	// Labels custom labels for grouping and filtering
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

var GlobalConfig *Config
var configFilePath string // Store the path of the loaded config file

// Load loads configuration file, supports using viper or yaml parsing
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("configuration file path cannot be empty")
	}
	configFilePath = path // Store for saving later

	ext := filepath.Ext(path)
	var cfg *Config
	var err error

	switch ext {
	case ".yaml", ".yml":
		// Try to load configuration using viper
		cfg, err = loadViperConfig(path)
		if err != nil {
			// If viper fails, fallback to original yaml parsing
			cfg, err = loadYAMLConfig(path)
		}
	default:
		return nil, fmt.Errorf("unsupported configuration file format: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	GlobalConfig = cfg
	setDefaults()

	return cfg, nil
}

// loadViperConfig loads configuration file using viper
func loadViperConfig(path string) (*Config, error) {
	v := viper.New()

	// Set configuration file path and name
	v.SetConfigFile(path)

	// Set environment variable prefix
	v.SetEnvPrefix("CILIKUBE")
	v.AutomaticEnv()

	// Set field name mapping so viper can correctly map fields
	v.RegisterAlias("server.activeCluster", "server.activeClusterID")

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read configuration file %s: %w", path, err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("viper failed to parse configuration file: %w", err)
	}

	return cfg, nil
}

func loadYAMLConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read configuration file %s: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML configuration file: %w", err)
	}

	return cfg, nil
}

// SaveGlobalConfig saves the current GlobalConfig to its original loading path
func SaveGlobalConfig() error {
	if GlobalConfig == nil {
		return fmt.Errorf("global configuration not yet loaded, cannot save")
	}
	if configFilePath == "" {
		return fmt.Errorf("configuration file path unknown, cannot save")
	}

	data, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		return fmt.Errorf("failed to serialize configuration to YAML: %w", err)
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp(filepath.Dir(configFilePath), filepath.Base(configFilePath)+".tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary configuration file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write temporary configuration file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary configuration file: %w", err)
	}

	// Replace the original file with the temporary file
	if err := os.Rename(tempFile.Name(), configFilePath); err != nil {
		return fmt.Errorf("failed to replace original configuration file: %w", err)
	}

	return nil
}

func setDefaults() {
	if GlobalConfig.Server.Port == "" {
		GlobalConfig.Server.Port = "8080"
	}
	if GlobalConfig.Server.Mode == "" {
		GlobalConfig.Server.Mode = "debug"
	}
	if GlobalConfig.Server.ReadTimeout == 0 {
		GlobalConfig.Server.ReadTimeout = 30
	}
	if GlobalConfig.Server.WriteTimeout == 0 {
		GlobalConfig.Server.WriteTimeout = 30
	}
	// ... (other default value settings for database, jwt, installer, kubernetes remain unchanged) ...
	if GlobalConfig.Database.Enabled { // Fix: only set database default values when enabled
		// Set default database type if not specified
		if GlobalConfig.Database.Type == "" {
			GlobalConfig.Database.Type = "mysql" // Default to MySQL for backward compatibility
		}

		// Set defaults based on database type
		switch GlobalConfig.Database.Type {
		case "sqlite":
			if GlobalConfig.Database.Database == "" {
				GlobalConfig.Database.Database = "./data/cilikube.db"
			}
			// SQLite doesn't need host, port, username, password
		case "postgresql", "postgres":
			if GlobalConfig.Database.Host == "" {
				GlobalConfig.Database.Host = "localhost"
			}
			if GlobalConfig.Database.Port == 0 {
				GlobalConfig.Database.Port = 5432 // PostgreSQL default port
			}
			if GlobalConfig.Database.Username == "" {
				GlobalConfig.Database.Username = "postgres"
			}
			if GlobalConfig.Database.Password == "" {
				GlobalConfig.Database.Password = "cilikube-password-change-in-production"
			}
			if GlobalConfig.Database.Database == "" {
				GlobalConfig.Database.Database = "cilikube"
			}
		case "mysql":
		default:
			if GlobalConfig.Database.Host == "" {
				GlobalConfig.Database.Host = "localhost"
			}
			if GlobalConfig.Database.Port == 0 {
				GlobalConfig.Database.Port = 3306 // MySQL default port
			}
			if GlobalConfig.Database.Username == "" {
				GlobalConfig.Database.Username = "root"
			}
			if GlobalConfig.Database.Password == "" {
				GlobalConfig.Database.Password = "cilikube-password-change-in-production"
			}
			if GlobalConfig.Database.Database == "" {
				GlobalConfig.Database.Database = "cilikube"
			}
			if GlobalConfig.Database.Charset == "" {
				GlobalConfig.Database.Charset = "utf8mb4"
			}
		}
	}

	if GlobalConfig.JWT.SecretKey == "" {
		GlobalConfig.JWT.SecretKey = os.Getenv("JWT_SECRET")
		if GlobalConfig.JWT.SecretKey == "" {
			GlobalConfig.JWT.SecretKey = "cilikube-secret-key-change-in-production"
		}
	}
	if GlobalConfig.JWT.ExpireDuration == 0 {
		GlobalConfig.JWT.ExpireDuration = 24 * time.Hour
	}
	if GlobalConfig.JWT.Issuer == "" {
		GlobalConfig.JWT.Issuer = "cilikube"
	}
	if GlobalConfig.Installer.MinikubeDriver == "" {
		GlobalConfig.Installer.MinikubeDriver = "docker"
	}
	if GlobalConfig.Installer.DownloadDir == "" {
		GlobalConfig.Installer.DownloadDir = "."
	}
	if GlobalConfig.Kubernetes.Kubeconfig == "" || GlobalConfig.Kubernetes.Kubeconfig == "default" {
		if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
			GlobalConfig.Kubernetes.Kubeconfig = kubeconfigEnv
		} else {
			home, err := os.UserHomeDir()
			if err == nil {
				GlobalConfig.Kubernetes.Kubeconfig = filepath.Join(home, ".kube", "config")
			} else {
				GlobalConfig.Kubernetes.Kubeconfig = ""
			}
		}
	}

	// Automatically generate UUID for clusters without ID
	configChanged := false
	var firstActiveClusterID string

	for i := range GlobalConfig.Clusters {
		if GlobalConfig.Clusters[i].ID == "" {
			GlobalConfig.Clusters[i].ID = uuid.New().String()
			configChanged = true
		}

		// Record the ID of the first active cluster for setting default active cluster
		if GlobalConfig.Clusters[i].IsActive && firstActiveClusterID == "" {
			firstActiveClusterID = GlobalConfig.Clusters[i].ID
		}
	}

	// If no active cluster ID is set, use the first active cluster's ID
	if GlobalConfig.Server.ActiveClusterID == "" && firstActiveClusterID != "" {
		GlobalConfig.Server.ActiveClusterID = firstActiveClusterID
		configChanged = true
	}

	// Set storage configuration defaults
	setStorageDefaults()

	// Set security configuration defaults
	setSecurityDefaults()

	// If new ID was generated or active cluster was updated, save configuration file
	if configChanged {
		_ = SaveGlobalConfig() // Ignore errors as this is optional
	}
}

func (c *Config) GetDSN() string {
	if !c.Database.Enabled {
		return "" // If database is not enabled, return empty DSN
	}

	// Support different database types
	switch c.Database.Type {
	case "sqlite":
		return c.Database.Database // For SQLite, database field contains the file path
	case "postgresql", "postgres":
		// PostgreSQL DSN format
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Database.Host,
			c.Database.Port,
			c.Database.Username,
			c.Database.Password,
			c.Database.Database)
	case "mysql", "":
		// Default to MySQL format for backward compatibility
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			c.Database.Username,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Database,
			c.Database.Charset)
	default:
		// Default to MySQL format
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			c.Database.Username,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Database,
			c.Database.Charset)
	}
}

// GetClusterByID gets cluster information by ID
func (c *Config) GetClusterByID(id string) *ClusterInfo {
	for i := range c.Clusters {
		if c.Clusters[i].ID == id {
			return &c.Clusters[i]
		}
	}
	return nil
}

// GetClusterByName gets cluster information by name (backward compatibility)
func (c *Config) GetClusterByName(name string) *ClusterInfo {
	for i := range c.Clusters {
		if c.Clusters[i].Name == name {
			return &c.Clusters[i]
		}
	}
	return nil
}

// GetClusterIDByName gets cluster ID by name (backward compatibility)
func (c *Config) GetClusterIDByName(name string) string {
	cluster := c.GetClusterByName(name)
	if cluster != nil {
		return cluster.ID
	}
	return ""
}

// setStorageDefaults sets default values for storage configuration
func setStorageDefaults() {
	// Logic for automatically selecting storage type
	if GlobalConfig.Storage.Type == "" {
		GlobalConfig.Storage.Type = DetermineStorageType(&GlobalConfig.Storage)
	}

	// If no storage database configuration is specified, use global database configuration
	if GlobalConfig.Storage.Database == nil && GlobalConfig.Storage.Type == "database" {
		GlobalConfig.Storage.Database = &GlobalConfig.Database
	}

	// Set OAuth default configuration
	if GlobalConfig.OAuth.GitHub.RedirectURL == "" {
		GlobalConfig.OAuth.GitHub.RedirectURL = "http://localhost:8080/api/v1/auth/oauth/callback"
	}
}

// DetermineStorageType automatically determines storage type based on configuration
func DetermineStorageType(config *StorageConfig) string {
	// If type is explicitly specified, use the specified type
	if config.Type != "" {
		return config.Type
	}

	// Check database settings in storage configuration
	var dbConfig *DatabaseConfig
	if config.Database != nil {
		dbConfig = config.Database
	} else {
		// If there are no database settings in storage configuration, check global database configuration
		dbConfig = &GlobalConfig.Database
	}

	// If there is no database configuration or database is not enabled, use memory storage by default
	if dbConfig == nil || !dbConfig.Enabled {
		return "memory"
	}

	// For SQLite, only need to check the database field (file path)
	if dbConfig.Type == "sqlite" {
		if dbConfig.Database == "" {
			return "memory"
		}
		return "database"
	}

	// For other database types (MySQL, PostgreSQL), need to check Host and Database
	if dbConfig.Host == "" || dbConfig.Database == "" {
		return "memory"
	}

	// Have complete database configuration, use database storage
	return "database"
}

// GetStorageType returns the determined storage type
func (c *Config) GetStorageType() string {
	return DetermineStorageType(&c.Storage)
}

// GetStorageDSN returns the DSN for storage database connection
func (c *Config) GetStorageDSN() string {
	if c.Storage.Type != "database" {
		return ""
	}

	dbConfig := c.Storage.Database
	if dbConfig == nil {
		dbConfig = &c.Database
	}

	if !dbConfig.Enabled {
		return ""
	}

	// Support different database types for storage
	switch dbConfig.Type {
	case "sqlite":
		return dbConfig.Database // For SQLite, database field contains the file path
	case "postgresql", "postgres":
		// PostgreSQL DSN format
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Database)
	case "mysql", "":
		// Default to MySQL format for backward compatibility
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
			dbConfig.Charset)
	default:
		// Default to MySQL format
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
			dbConfig.Charset)
	}
}

// setSecurityDefaults sets default values for security configuration
func setSecurityDefaults() {
	// Password policy defaults
	if GlobalConfig.Security.Password.MinLength == 0 {
		GlobalConfig.Security.Password.MinLength = 8
	}
	// Default to requiring at least lowercase and numbers for basic security
	if !GlobalConfig.Security.Password.RequireUppercase &&
		!GlobalConfig.Security.Password.RequireLowercase &&
		!GlobalConfig.Security.Password.RequireNumbers &&
		!GlobalConfig.Security.Password.RequireSymbols {
		GlobalConfig.Security.Password.RequireLowercase = true
		GlobalConfig.Security.Password.RequireNumbers = true
	}

	// Account lockout defaults
	if GlobalConfig.Security.AccountLock.MaxFailedAttempts == 0 {
		GlobalConfig.Security.AccountLock.MaxFailedAttempts = 5
	}
	if GlobalConfig.Security.AccountLock.LockoutDuration == 0 {
		GlobalConfig.Security.AccountLock.LockoutDuration = 15 * time.Minute
	}
	if GlobalConfig.Security.AccountLock.ResetAfter == 0 {
		GlobalConfig.Security.AccountLock.ResetAfter = 1 * time.Hour
	}

	// Session management defaults
	if GlobalConfig.Security.Session.MaxConcurrentSessions == 0 {
		GlobalConfig.Security.Session.MaxConcurrentSessions = 3
	}
	if GlobalConfig.Security.Session.IdleTimeout == 0 {
		GlobalConfig.Security.Session.IdleTimeout = 30 * time.Minute
	}
	if GlobalConfig.Security.Session.AbsoluteTimeout == 0 {
		GlobalConfig.Security.Session.AbsoluteTimeout = 8 * time.Hour
	}

	// Rate limiting defaults
	if GlobalConfig.Security.RateLimit.LoginAttempts == 0 {
		GlobalConfig.Security.RateLimit.LoginAttempts = 10
	}
	if GlobalConfig.Security.RateLimit.LoginWindow == 0 {
		GlobalConfig.Security.RateLimit.LoginWindow = 15 * time.Minute
	}
	if GlobalConfig.Security.RateLimit.APIRequests == 0 {
		GlobalConfig.Security.RateLimit.APIRequests = 1000
	}
	if GlobalConfig.Security.RateLimit.APIWindow == 0 {
		GlobalConfig.Security.RateLimit.APIWindow = 1 * time.Hour
	}
	if GlobalConfig.Security.RateLimit.BurstSize == 0 {
		GlobalConfig.Security.RateLimit.BurstSize = 50
	}
}
