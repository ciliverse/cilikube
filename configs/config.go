package configs

import (
	"os"
)

type Config struct {
	KubeConfig string
}

func LoadConfig() *Config {
	return &Config{
		KubeConfig: os.Getenv("KUBECONFIG"),
	}
}
