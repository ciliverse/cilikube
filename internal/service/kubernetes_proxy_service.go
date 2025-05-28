package service

import (
	"k8s.io/client-go/rest"
)

type ProxyService struct {
	restConfig *rest.Config
}

func NewProxyService(restConfig *rest.Config) *ProxyService {
	return &ProxyService{
		restConfig: restConfig,
	}
}

func (s *ProxyService) GetConfig() *rest.Config {
	return s.restConfig
}
