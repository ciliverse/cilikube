package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client 结构体封装了与 Kubernetes 集群交互所需的 clientset 和 rest.Config。
type Client struct {
	Clientset kubernetes.Interface
	Config    *rest.Config
}

// NewClient 通过指定的 kubeconfig 文件路径创建一个新的 Kubernetes 客户端。
// 它能够处理 "in-cluster"（在集群内部署）、"default"（默认路径）以及明确的件路径。
func NewClient(kubeconfig string) (*Client, error) {
	var config *rest.Config
	var err error

	// 当应用部署在 Kubernetes 集群内部时，使用 "in-cluster" 配置。
	if kubeconfig == "in-cluster" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("加载 in-cluster 配置失败: %w", err)
		}
	} else {
		// 当 kubeconfig 为空或 "default" 时，使用标准的用户主目录下的 .kube/config 文件。
		if kubeconfig == "default" || kubeconfig == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("获取用户主目录失败: %w", err)
			}
			kubeconfig = filepath.Join(homeDir, ".kube", "config")
		}

		// 检查 kubeconfig 文件是否存在。
		if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig 文件不存在: %s", kubeconfig)
		}

		// 从指定的 kubeconfig 文件路径构建配置。
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("从 kubeconfig 文件 '%s' 构建配置失败: %w", kubeconfig, err)
		}
	}

	// 使用生成的配置创建 clientset。
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes clientset 失败: %w", err)
	}

	return &Client{
		Clientset: clientset,
		Config:    config,
	}, nil
}

// NewClientFromContent 通过内存中的 kubeconfig 文件内容创建一个新的 Kubernetes 客户端。
// 这是实现动态集群管理（例如，从前端API接收kubeconfig）的关键函数。
func NewClientFromContent(kubeconfigData []byte) (*Client, error) {
	if len(kubeconfigData) == 0 {
		return nil, fmt.Errorf("kubeconfig 内容不能为空")
	}

	// 从字节切片创建 clientcmd 配置。
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
	if err != nil {
		return nil, fmt.Errorf("从字节内容创建客户端配置失败: %w", err)
	}

	// 从 clientcmd 配置中获取 REST 客户端配置。
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("从客户端配置中获取 REST 配置失败: %w", err)
	}

	// 使用 REST 配置创建 clientset。
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes clientset 失败: %w", err)
	}

	return &Client{
		Clientset: clientset,
		Config:    restConfig,
	}, nil
}

// CheckConnection 对 Kubernetes API Server 执行一次轻量级的健康检查。
func (c *Client) CheckConnection() error {
	if c == nil || c.Clientset == nil {
		return fmt.Errorf("kubernetes 客户端未初始化")
	}
	// 使用 discovery 客户端获取服务器版本作为一种低成本的连接检查方式。
	_, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("检查 Kubernetes API Server 连接失败: %w", err)
	}
	return nil
}
