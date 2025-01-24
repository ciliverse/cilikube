package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset kubernetes.Interface
}

func NewClient(kubeconfig string) (*Client, error) {
	var config *rest.Config
	var err error

	// 优先使用in-cluster配置
	if kubeconfig == "" {
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	} else {
		// 处理kubeconfig路径
		if kubeconfig == "default" {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{Clientset: clientset}, nil
}

// 健康检查扩展
func (c *Client) CheckConnection() error {
	_, err := c.Clientset.Discovery().ServerVersion()
	return err
}
