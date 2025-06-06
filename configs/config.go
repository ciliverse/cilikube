package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	Port          string `yaml:"port" json:"port"`
	ReadTimeout   int    `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout  int    `yaml:"write_timeout" json:"write_timeout"`
	Mode          string `yaml:"mode" json:"mode"` // debug, release
	ActiveCluster string `yaml:"activeCluster" json:"activeCluster"`
	// [NEW] 添加用于加密 kubeconfig 的密钥
	EncryptionKey string `yaml:"encryptionKey" json:"encryptionKey"`
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
	Username string `yaml:"username" json:"username"` // 确保这里是 username
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"` // 确保这里是 database
	Charset  string `yaml:"charset" json:"charset"`
}

type JWTConfig struct {
	SecretKey      string        `yaml:"secret_key" json:"secret_key"`
	ExpireDuration time.Duration `yaml:"expire_duration" json:"expire_duration"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
}

type ClusterInfo struct {
	Name string `yaml:"name" json:"name"` // 集群的唯一标识名称
	// ConfigPath 可以是 kubeconfig 文件的绝对路径，或者是 "in-cluster"
	ConfigPath string `yaml:"config_path" json:"config_path"`
	IsActive   bool   `yaml:"is_active" json:"is_active"` // 标记此集群配置是否启用
}

var GlobalConfig *Config
var configFilePath string // Store the path of the loaded config file

// Load 加载配置文件
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("配置文件路径不能为空")
	}
	configFilePath = path // Store for saving later

	ext := filepath.Ext(path)
	var cfg *Config
	var err error

	switch ext {
	case ".yaml", ".yml":
		cfg, err = loadYAMLConfig(path)
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	GlobalConfig = cfg
	setDefaults()

	return cfg, nil
}

func loadYAMLConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析 YAML 配置文件失败: %w", err)
	}

	return cfg, nil
}

// SaveGlobalConfig 将当前的 GlobalConfig 保存到其原始加载路径
func SaveGlobalConfig() error {
	if GlobalConfig == nil {
		return fmt.Errorf("全局配置尚未加载, 无法保存")
	}
	if configFilePath == "" {
		return fmt.Errorf("配置文件路径未知, 无法保存")
	}

	data, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		return fmt.Errorf("序列化配置到 YAML 失败: %w", err)
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp(filepath.Dir(configFilePath), filepath.Base(configFilePath)+".tmp")
	if err != nil {
		return fmt.Errorf("创建临时配置文件失败: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("写入临时配置文件失败: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("关闭临时配置文件失败: %w", err)
	}

	// Replace the original file with the temporary file
	if err := os.Rename(tempFile.Name(), configFilePath); err != nil {
		return fmt.Errorf("替换原配置文件失败: %w", err)
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
	// ... (其他 database, jwt, installer, kubernetes 的默认值设置保持不变) ...
	if GlobalConfig.Database.Enabled { // 修正：只在 enabled 时设置数据库默认值
		if GlobalConfig.Database.Host == "" {
			GlobalConfig.Database.Host = "localhost"
		}
		if GlobalConfig.Database.Port == 0 {
			// 对于 MySQL 通常是 3306，PostgreSQL 是 5432。这里以 MySQL 为例。
			GlobalConfig.Database.Port = 3306
		}
		if GlobalConfig.Database.Username == "" { // 对应 DatabaseConfig 中的 Username
			GlobalConfig.Database.Username = "root"
		}
		if GlobalConfig.Database.Password == "" {
			GlobalConfig.Database.Password = "cilikube-password-change-in-production"
		}
		if GlobalConfig.Database.Database == "" { // 对应 DatabaseConfig 中的 Database
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
}

func (c *Config) GetDSN() string {
	if !c.Database.Enabled {
		return "" // 如果数据库未启用，返回空 DSN
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.Database.Username, // 确保这里是 Username
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database, // 确保这里是 Database
		c.Database.Charset)
}
