package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings" // 引入 strings 包用于 TrimSpace

	"gopkg.in/yaml.v3"
)

// ServerConfig 定义服务器相关配置
type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	// ActiveCluster 字段可以保留，ClientManager 可以用它来决定初始的默认集群，
	// 或者你可以依赖下面 IndividualClusterConfig 中的 IsDefault 字段。
	// 如果两者都用，需要明确优先级。这里暂时注释掉，优先使用 IsDefault。
	// ActiveCluster string `yaml:"activeCluster"`
}

// IndividualClusterConfig 定义了单个被管理 Kubernetes 集群的配置信息
type IndividualClusterConfig struct {
	Name       string `yaml:"name"`       // 集群的唯一名称，例如 "prod-cluster", "dev-local"
	Kubeconfig string `yaml:"kubeconfig"` // Kubeconfig 文件路径。"default" 表示 ~/.kube/config, "" 表示 in-cluster
	IsDefault  bool   `yaml:"isDefault"`  // 可选: 如果为 true，则此集群作为默认集群
}

// KubernetesMultiClusterConfig 替代了旧的 KubernetesConfig，用于管理多个集群的配置
type KubernetesMultiClusterConfig struct {
	Clusters []IndividualClusterConfig `yaml:"clusters"` // 管理的集群列表
	// 你也可以在这里定义一个顶层的 defaultClusterName 来指定默认集群的名称
	// DefaultClusterName string `yaml:"defaultClusterName"`
}

// InstallerConfig 定义了安装器 (如 Minikube) 的相关配置 (保持不变)
type InstallerConfig struct {
	MinikubePath   string `yaml:"minikubePath"`
	MinikubeDriver string `yaml:"minikubeDriver"`
	DownloadDir    string `yaml:"downloadDir"`
}

// Config 是应用的总配置结构体
type Config struct {
	Server     ServerConfig                 `yaml:"server"`
	Kubernetes KubernetesMultiClusterConfig `yaml:"kubernetes"` // 这里从旧的 KubernetesConfig 更新为 KubernetesMultiClusterConfig
	Installer  InstallerConfig              `yaml:"installer"`
}

// GlobalAppConfig 用于存储加载后的全局配置实例
// 其他包可以通过 configs.GlobalAppConfig 访问配置
var GlobalAppConfig *Config

// Load 从指定的 YAML 文件路径加载配置
func Load(path string) (*Config, error) {
	// 默认配置文件路径
	if path == "" {
		// 确保这个默认路径与你的项目结构和部署方式一致
		// 例如，可以尝试 "./configs/config.yaml" 或相对于可执行文件的路径
		path = "configs/config.yaml"
		log.Printf("信息: 未提供配置文件路径，使用默认路径: %s", path)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件 '%s' 不存在", path)
	}

	// 读取配置文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 '%s': %w", path, err)
	}

	// 解析 YAML 数据到 Config 结构体
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件 '%s' 失败: %w", path, err)
	}

	// --- 设置和处理默认值 ---

	// Server 配置默认值
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	// Installer 配置默认值
	if cfg.Installer.MinikubeDriver == "" {
		cfg.Installer.MinikubeDriver = "docker"
	}
	if cfg.Installer.DownloadDir == "" {
		// 考虑使用更合适的用户特定目录作为默认，例如 os.UserCacheDir()
		cfg.Installer.DownloadDir = "./cilikube_downloads" // 示例：当前目录下的子目录
		log.Printf("信息: Installer.DownloadDir 未设置，默认为: %s", cfg.Installer.DownloadDir)
	}

	// Kubernetes 多集群配置处理
	if len(cfg.Kubernetes.Clusters) == 0 {
		log.Println("警告: 在 'kubernetes.clusters' 下未配置任何集群。应用可能无法连接到任何 Kubernetes 环境。")
		// 如果你想支持旧的单 kubeconfig 格式作为回退，这里的逻辑会更复杂。
		// 为了清晰地转向多集群，建议要求用户迁移到 `clusters` 数组格式。
	}

	defaultClusterExplicitlySet := false
	var defaultClusterName string

	for i := range cfg.Kubernetes.Clusters {
		// 获取集群配置的指针，以便直接修改（例如解析后的kubeconfig路径）
		clusterCfg := &cfg.Kubernetes.Clusters[i]

		// 校验集群名称是否存在
		if clusterCfg.Name == "" {
			// 对于 ClientManager 来说，集群名称是必要的标识符
			return nil, fmt.Errorf("错误: 'kubernetes.clusters' 定义中的第 %d 个集群缺少 'name' 字段", i+1)
		}
		clusterCfg.Name = strings.TrimSpace(clusterCfg.Name)

		// 清理并解析 kubeconfig 路径
		clusterCfg.Kubeconfig = strings.TrimSpace(clusterCfg.Kubeconfig)
		if clusterCfg.Kubeconfig == "default" {
			// 优先使用 KUBECONFIG 环境变量
			if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
				clusterCfg.Kubeconfig = kubeconfigEnv
				log.Printf("信息: 集群 '%s' 的 kubeconfig 为 'default'，使用 KUBECONFIG 环境变量: %s\n", clusterCfg.Name, clusterCfg.Kubeconfig)
			} else {
				// 否则，使用用户主目录下的 .kube/config
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return nil, fmt.Errorf("无法获取用户主目录以解析集群 '%s' 的 'default' kubeconfig: %w", clusterCfg.Name, err)
				}
				clusterCfg.Kubeconfig = filepath.Join(homeDir, ".kube", "config")
				log.Printf("信息: 集群 '%s' 的 kubeconfig 为 'default'，解析为默认路径: %s\n", clusterCfg.Name, clusterCfg.Kubeconfig)
			}
		} else if clusterCfg.Kubeconfig == "" {
			// 空字符串表示尝试 in-cluster 配置，这将在 NewClient 中处理
			log.Printf("信息: 集群 '%s' 的 kubeconfig 为空字符串，将尝试使用 in-cluster 配置。\n", clusterCfg.Name)
		} else {
			// 用户指定了具体路径
			// 可选: 检查路径有效性或转换为绝对路径
			// if !filepath.IsAbs(clusterCfg.Kubeconfig) { ... }
			log.Printf("信息: 集群 '%s' 使用指定的 kubeconfig 路径: %s\n", clusterCfg.Name, clusterCfg.Kubeconfig)
		}

		// 处理 IsDefault 标志
		if clusterCfg.IsDefault {
			if defaultClusterExplicitlySet {
				log.Printf("警告: 找到多个 'isDefault: true' 的集群。将使用第一个被标记为默认的集群 ('%s')。集群 '%s' 的 IsDefault 标志将被忽略。", defaultClusterName, clusterCfg.Name)
			} else {
				defaultClusterExplicitlySet = true
				defaultClusterName = clusterCfg.Name
				log.Printf("信息: 集群 '%s' 被配置为默认集群。", clusterCfg.Name)
			}
		}
	}

	// 如果在循环后没有显式设置的默认集群，但 ServerConfig.ActiveCluster (如果启用) 或 KubernetesMultiClusterConfig.DefaultClusterName (如果启用) 被设置，
	// 可以在这里或 ClientManager 初始化时使用它们来确定默认集群。
	// 例如，如果启用了 ServerConfig.ActiveCluster:
	// if !defaultClusterExplicitlySet && cfg.Server.ActiveCluster != "" {
	//     foundActive := false
	//     for _, cluster := range cfg.Kubernetes.Clusters {
	//         if cluster.Name == cfg.Server.ActiveCluster {
	//             log.Printf("信息: 使用 server.activeCluster ('%s') 作为默认集群。", cfg.Server.ActiveCluster)
	//             // ClientManager 在初始化时会用到这个信息
	//             defaultClusterName = cfg.Server.ActiveCluster // 更新 defaultClusterName
	//             defaultClusterExplicitlySet = true
	//             foundActive = true
	//             break
	//         }
	//     }
	//     if !foundActive {
	//         log.Printf("警告: server.activeCluster 指定的集群 '%s' 未在 'kubernetes.clusters' 列表中找到。", cfg.Server.ActiveCluster)
	//     }
	// }

	log.Printf("配置文件 '%s' 加载并初步处理完成。共定义了 %d 个集群。", path, len(cfg.Kubernetes.Clusters))

	// 将加载并处理后的配置存储到全局变量
	GlobalAppConfig = &cfg
	return &cfg, nil
}
