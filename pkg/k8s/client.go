package k8s

import (
	"fmt" // Import fmt for errors
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client struct now holds both Clientset and the Config
type Client struct {
	Clientset kubernetes.Interface
	Config    *rest.Config // <-- 添加 Config 字段来存储 rest.Config
}

// NewClient creates a new Kubernetes client instance.
// It attempts in-cluster config first if kubeconfig is empty,
// otherwise uses the provided kubeconfig path.
func NewClient(kubeconfig string) (*Client, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		// Attempt in-cluster configuration (typically for production)
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("加载 in-cluster 配置失败: %w", err)
		}
		fmt.Println("使用 in-cluster Kubernetes 配置。") // Log which config is used
	} else {
		// Use out-of-cluster configuration (typically for development)
		// Handle 'default' keyword to find config in home directory
		if kubeconfig == "default" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("无法获取用户主目录: %w", err)
			}
			kubeconfig = filepath.Join(homeDir, ".kube", "config")
			fmt.Printf("使用默认 kubeconfig 路径: %s\n", kubeconfig) // Log which config is used
		} else {
			fmt.Printf("使用指定的 kubeconfig 路径: %s\n", kubeconfig) // Log which config is used
		}

		// Check if the resolved kubeconfig file exists
		if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig 文件不存在: %s", kubeconfig)
		}

		// Build configuration from flags
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("从 kubeconfig 构建配置失败 '%s': %w", kubeconfig, err)
		}
	}

	// Create the clientset using the configuration
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes clientset 失败: %w", err)
	}

	// Return the Client struct containing BOTH clientset and config
	return &Client{
		Clientset: clientset,
		Config:    config, // <-- 将加载的 config 存储在结构体中
	}, nil
}

// CheckConnection performs a basic health check against the Kubernetes API server.
func (c *Client) CheckConnection() error {
	// Ensure clientset is not nil before using
	if c == nil || c.Clientset == nil {
		return fmt.Errorf("kubernetes client 未初始化")
	}
	_, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("检查 Kubernetes 连接失败: %w", err)
	}
	return nil
}
