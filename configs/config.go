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
	JWT        JWTConfig        `yaml:"jwt" json:"jwt"`
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
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"` // Ensure this is username
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"` // Ensure this is database
	Charset  string `yaml:"charset" json:"charset"`
}

type JWTConfig struct {
	SecretKey      string        `yaml:"secret_key" json:"secret_key"`
	ExpireDuration time.Duration `yaml:"expire_duration" json:"expire_duration"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
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
		if GlobalConfig.Database.Host == "" {
			GlobalConfig.Database.Host = "localhost"
		}
		if GlobalConfig.Database.Port == 0 {
			// For MySQL it's usually 3306, PostgreSQL is 5432. Here we use MySQL as example.
			GlobalConfig.Database.Port = 3306
		}
		if GlobalConfig.Database.Username == "" { // Corresponds to Username in DatabaseConfig
			GlobalConfig.Database.Username = "root"
		}
		if GlobalConfig.Database.Password == "" {
			GlobalConfig.Database.Password = "cilikube-password-change-in-production"
		}
		if GlobalConfig.Database.Database == "" { // Corresponds to Database in DatabaseConfig
			GlobalConfig.Database.Database = "cilikube"
		}
		if GlobalConfig.Database.Charset == "" {
			GlobalConfig.Database.Charset = "utf8mb4"
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

	// If new ID was generated or active cluster was updated, save configuration file
	if configChanged {
		_ = SaveGlobalConfig() // Ignore errors as this is optional
	}
}

func (c *Config) GetDSN() string {
	if !c.Database.Enabled {
		return "" // If database is not enabled, return empty DSN
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.Database.Username, // Ensure this is Username
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database, // Ensure this is Database
		c.Database.Charset)
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
