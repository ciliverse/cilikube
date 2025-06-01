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
	Name string `yaml:"name" json:"name"` // 集群的唯一标识名称，将用于API路径或参数
	// ConfigPath 可以是 kubeconfig 文件的绝对路径，或者是 "in-cluster"（如果管理平台本身部署在目标集群内并希望使用服务账户）
	ConfigPath string `yaml:"config_path" json:"config_path"`
	IsActive   bool   `yaml:"is_active" json:"is_active"` // 标记此集群配置是否启用
	// 可以考虑添加的字段：
	// DisplayName string            `yaml:"display_name" json:"display_name"` // 用于UI显示的名称
	// Description string            `yaml:"description" json:"description"`   // 集群描述
	// Context     string            `yaml:"context,omitempty" json:"context"` // 如果 kubeconfig 文件包含多个 context，指定使用哪一个
	// ReadOnly    bool              `yaml:"read_only,omitempty" json:"read_only"` // 是否只读集群
}

var GlobalConfig *Config

// Load 加载配置文件
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("配置文件路径不能为空")
	}

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
	setDefaults() // 保持原有的默认值设定逻辑

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
		if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
			GlobalConfig.Kubernetes.Kubeconfig = kubeconfigEnv
		} else {
			home, err := os.UserHomeDir()
			if err == nil { // 只有成功获取到主目录才设置
				GlobalConfig.Kubernetes.Kubeconfig = filepath.Join(home, ".kube", "config")
			} else {
				// 如果无法获取用户主目录，且环境变量也未设置，则此路径可能无效或为空。
				// 此时应用应能优雅处理，例如日志警告，并依赖 `Clusters` 配置。
				GlobalConfig.Kubernetes.Kubeconfig = "" // 或者一个明确的标记，表示未配置
			}
		}
	}

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
			// 确保为数据库密码提供一个安全的默认值或强制用户配置
			GlobalConfig.Database.Password = "cilikube-password-change-in-production"
		}
		if GlobalConfig.Database.Charset == "" {
			GlobalConfig.Database.Charset = "utf8mb4"
		}
	}

	// for i := range GlobalConfig.Clusters {
	//    if GlobalConfig.Clusters[i].Name == "" {
	//        log.Printf("警告: 第 %d 个集群配置缺少名称，这可能导致问题。", i+1)
	//    }
	//    // 默认 IsActive 为 true，除非显式设置为 false？
	//    // if !isSet(GlobalConfig.Clusters[i].IsActive) { //伪代码: isSet 需要具体实现
	//    //    GlobalConfig.Clusters[i].IsActive = true
	//    // }
	// }
}

func (c *Config) GetDSN() string {
	// 保持原有的 DSN 生成逻辑
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset)
}
