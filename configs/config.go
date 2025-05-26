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
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
	Charset  string `yaml:"charset" json:"charset"`
}

type JWTConfig struct {
	SecretKey      string        `yaml:"secret_key" json:"secret_key"`
	ExpireDuration time.Duration `yaml:"expire_duration" json:"expire_duration"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
}

type ClusterInfo struct {
	Name       string `yaml:"name" json:"name"`
	ConfigPath string `yaml:"config_path" json:"config_path"`
	IsActive   bool   `yaml:"is_active" json:"is_active"`
}

var GlobalConfig *Config

// Load 加载配置文件
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("配置文件路径不能为空")
	}

	// 检查文件扩展名
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

func setDefaults() {
	if GlobalConfig.Server.Port == "" {
		GlobalConfig.Server.Port = "8080"
	}
	if GlobalConfig.Server.Mode == "" {
		GlobalConfig.Server.Mode = "debug"
	}
	if GlobalConfig.Server.ReadTimeout == 0 {
		GlobalConfig.Server.ReadTimeout = 30 // 默认 30 秒
	}
	if GlobalConfig.Server.WriteTimeout == 0 {
		GlobalConfig.Server.WriteTimeout = 30 // 默认 30 秒
	}
	if GlobalConfig.Database.Host == "" {
		GlobalConfig.Database.Host = "localhost"
	}
	if GlobalConfig.Database.Port == 0 {
		GlobalConfig.Database.Port = 3306
	}
	if GlobalConfig.Database.Username == "" {
		GlobalConfig.Database.Username = "root"
	}
	if GlobalConfig.Database.Charset == "" {
		GlobalConfig.Database.Charset = "utf8mb4"
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
		if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
			GlobalConfig.Kubernetes.Kubeconfig = kubeconfig
		} else {
			GlobalConfig.Kubernetes.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
	}
	// 数据库默认值
	if GlobalConfig.Database.Enabled {
		if GlobalConfig.Database.Host == "" {
			GlobalConfig.Database.Host = "localhost"
		}
		if GlobalConfig.Database.Port == 0 {
			GlobalConfig.Database.Port = 3306
		}
		if GlobalConfig.Database.Username == "" {
			GlobalConfig.Database.Username = "root"
		}
		if GlobalConfig.Database.Database == "" {
			GlobalConfig.Database.Database = "cilikube"
		}
		if GlobalConfig.Database.Password == "" {
			GlobalConfig.Database.Password = "cilikube-password-change-in-production"
		}
		if GlobalConfig.Database.Charset == "" {
			GlobalConfig.Database.Charset = "utf8mb4"
		}
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset)
}
