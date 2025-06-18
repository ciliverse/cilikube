package service

import (
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// ResourceServiceFactory 资源服务工厂
type ResourceServiceFactory struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewResourceServiceFactory 创建资源服务工厂
func NewResourceServiceFactory() *ResourceServiceFactory {
	return &ResourceServiceFactory{
		services: make(map[string]interface{}),
	}
}

// RegisterService 注册资源服务
func (f *ResourceServiceFactory) RegisterService(name string, service interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.services[name] = service
}

// GetService 获取资源服务
func (f *ResourceServiceFactory) GetService(name string) interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.services[name]
}

// InitializeDefaultServices 初始化默认服务
func (f *ResourceServiceFactory) InitializeDefaultServices() {
	// 注册所有资源服务
	f.RegisterService("services", NewBaseResourceService[*v1.Service]())
	f.RegisterService("pods", NewBaseResourceService[*v1.Pod]())
	f.RegisterService("deployments", NewBaseResourceService[*appsv1.Deployment]())
	f.RegisterService("daemonsets", NewBaseResourceService[*appsv1.DaemonSet]())
	f.RegisterService("ingresses", NewBaseResourceService[*networkingv1.Ingress]())
	f.RegisterService("configmaps", NewBaseResourceService[*v1.ConfigMap]())
	f.RegisterService("secrets", NewBaseResourceService[*v1.Secret]())
	f.RegisterService("persistentvolumeclaims", NewBaseResourceService[*v1.PersistentVolumeClaim]())
	f.RegisterService("persistentvolumes", NewBaseResourceService[*v1.PersistentVolume]())
	f.RegisterService("statefulsets", NewBaseResourceService[*appsv1.StatefulSet]())
	f.RegisterService("namespaces", NewBaseResourceService[*v1.Namespace]())
}
