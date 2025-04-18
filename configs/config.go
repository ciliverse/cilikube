package configs

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type KubernetesConfig struct {
	Kubeconfig string `yaml:"kubeconfig"`
}

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Kubernetes KubernetesConfig `yaml:"kubernetes"`
	Installer  InstallerConfig  `yaml:"installer"`
}

// InstallerConfig holds settings for the Minikube installer.
// <-- Add this struct definition -->
type InstallerConfig struct {
	// Base path where minikube executable might be found or should be (simulated) installed.
	// Leave empty to rely on PATH lookup.
	// Minikube可执行文件可能被找到或应该被（模拟）安装的基本路径。
	// 留空以依赖 PATH 查找。
	MinikubePath string `yaml:"minikubePath"`
	// The driver to use for minikube start (e.g., "docker", "hyperkit", "virtualbox").
	// minikube start 使用的驱动程序（例如 "docker", "hyperkit", "virtualbox"）。
	MinikubeDriver string `yaml:"minikubeDriver"`
	// Directory for temporary downloads. Defaults to current directory if empty.
	// 用于临时下载的目录。如果为空，则默认为当前目录。
	DownloadDir string `yaml:"downloadDir"`
}


func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
// --- Set Defaults ---
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080" // Default port
	}
	if cfg.Installer.MinikubeDriver == "" {
		cfg.Installer.MinikubeDriver = "docker" // Default driver
	}
	if cfg.Installer.DownloadDir == "" {
		cfg.Installer.DownloadDir = "." // Default to current directory
	}


	// 处理默认kubeconfig路径
	if cfg.Kubernetes.Kubeconfig == "default" {
		cfg.Kubernetes.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}


	return &cfg, nil
}
