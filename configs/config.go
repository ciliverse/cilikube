package configs

import (
	"fmt"
	"log"
	"os"

	// Keep for default kubeconfig path logic if needed elsewhere
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// KubernetesConfig is now DEPRECATED for single cluster, will be replaced by Clusters.
// We keep it temporarily if other parts of your app might still use it for non-multi-cluster features,
// but ideally, it should be removed. For now, we'll ignore it for multi-cluster setup.
// type KubernetesConfig struct {
// 	Kubeconfig string `yaml:"kubeconfig"` // Path for a single, default cluster
// }

type InstallerConfig struct {
	MinikubePath   string `yaml:"minikubePath"`
	MinikubeDriver string `yaml:"minikubeDriver"`
	DownloadDir    string `yaml:"downloadDir"`
}

// ClusterDefinition holds the configuration for a single Kubernetes cluster.
type ClusterDefinition struct {
	ID         string `yaml:"id"`         // Unique identifier for the cluster
	Name       string `yaml:"name"`       // Display name for the cluster
	Kubeconfig string `yaml:"kubeconfig"` // Kubeconfig content as a string
	// Add IsDefault bool `yaml:"isDefault"` if you want a default cluster concept
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	// Kubernetes KubernetesConfig    `yaml:"kubernetes"` // Deprecated or for single-cluster fallback
	Installer InstallerConfig     `yaml:"installer"`
	Clusters  []ClusterDefinition `yaml:"clusters"` // List of managed Kubernetes clusters
}

// GlobalAppConfig holds the loaded application configuration.
// Renamed from cfg to avoid conflict in Load function.
var GlobalAppConfig *Config

// Load parses the application configuration file.
func Load(path string) (*Config, error) {
	// 默认配置文件路径
	if path == "" {
		path = "configs/config.yaml" // Default path within the project
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %w", path, err)
	}

	// 解析配置文件
	var tempCfg Config
	if err := yaml.Unmarshal(data, &tempCfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	if tempCfg.Server.Port == "" {
		tempCfg.Server.Port = "8080"
	}
	if tempCfg.Installer.MinikubeDriver == "" {
		tempCfg.Installer.MinikubeDriver = "docker"
	}
	if tempCfg.Installer.DownloadDir == "" {
		tempCfg.Installer.DownloadDir = "."
	}

	// 旧的单个 Kubeconfig 默认值逻辑 (如果保留 KubernetesConfig 字段)
	// if tempCfg.Kubernetes.Kubeconfig == "default" {
	// 	tempCfg.Kubernetes.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// }

	// Validate cluster configurations
	if len(tempCfg.Clusters) == 0 {
		log.Println("警告: 配置文件中未定义任何集群 (clusters 列表为空或未定义)。多集群功能可能受限。")
	} else {
		for i, cluster := range tempCfg.Clusters {
			if cluster.ID == "" {
				return nil, fmt.Errorf("集群配置错误: 第 %d 个集群缺少 'id'", i+1)
			}
			if cluster.Kubeconfig == "" {
				return nil, fmt.Errorf("集群配置错误: 集群 '%s' (ID: %s) 缺少 'kubeconfig' 内容", cluster.Name, cluster.ID)
			}
			log.Printf("已加载集群定义: ID=%s, Name=%s", cluster.ID, cluster.Name)
		}
	}
	GlobalAppConfig = &tempCfg // Assign to global variable
	log.Printf("成功加载配置文件: %s", path)
	return GlobalAppConfig, nil
}

// GetClusterDefinitionByID retrieves a specific cluster's definition by its ID.
func GetClusterDefinitionByID(clusterID string) (*ClusterDefinition, bool) {
	if GlobalAppConfig == nil {
		return nil, false
	}
	for _, c := range GlobalAppConfig.Clusters {
		if c.ID == clusterID {
			return &c, true
		}
	}
	return nil, false
}

// GetClusterDefinitionByName retrieves a specific cluster's definition by its name.
func GetClusterDefinitionByName(clusterName string) (*ClusterDefinition, bool) {
	if GlobalAppConfig == nil {
		return nil, false
	}
	for _, c := range GlobalAppConfig.Clusters {
		if c.Name == clusterName {
			return &c, true
		}
	}
	return nil, false
}

// GetDefaultClusterDefinition retrieves the default cluster's definition.
func GetDefaultClusterDefinition() (*ClusterDefinition, bool) {
	if GlobalAppConfig == nil {
		return nil, false
	}
	for _, c := range GlobalAppConfig.Clusters {
		if c.ID == "default" { // Assuming "default" is the ID for the default cluster
			return &c, true
		}
	}
	return nil, false
}

// GetAllClusterDefinitions retrieves all cluster definitions.
func GetAllClusterDefinitions() []ClusterDefinition {
	if GlobalAppConfig == nil {
		return nil
	}
	return GlobalAppConfig.Clusters
}
