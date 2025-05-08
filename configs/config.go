package configs

import (
	"fmt"
	"log"
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

type InstallerConfig struct {
	MinikubePath   string `yaml:"minikubePath"`
	MinikubeDriver string `yaml:"minikubeDriver"`
	DownloadDir    string `yaml:"downloadDir"`
}

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Kubernetes KubernetesConfig `yaml:"kubernetes"`
	Installer  InstallerConfig  `yaml:"installer"`
}

func Load(path string) (*Config, error) {
	// 默认配置文件路径
	if path == "" {
		path = "config.yaml"
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
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Installer.MinikubeDriver == "" {
		cfg.Installer.MinikubeDriver = "docker"
	}
	if cfg.Installer.DownloadDir == "" {
		cfg.Installer.DownloadDir = "."
	}
	if cfg.Kubernetes.Kubeconfig == "default" {
		// 通过环境变量 KUBECONFIG指定配置文件
		if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
			cfg.Kubernetes.Kubeconfig = kubeconfig
		} else {
			cfg.Kubernetes.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}
	}

	log.Printf("成功加载配置文件: %s", path)
	return &cfg, nil
}
