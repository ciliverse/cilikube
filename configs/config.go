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

	// 处理默认kubeconfig路径
	if cfg.Kubernetes.Kubeconfig == "default" {
		cfg.Kubernetes.Kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	return &cfg, nil
}
