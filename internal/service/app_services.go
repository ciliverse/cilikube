package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// AppServices 应用程序服务集合
type AppServices struct {
	// 集群服务
	ClusterService *ClusterService

	// 资源服务
	ServiceService     *BaseResourceService[*corev1.Service]
	PodService         *BaseResourceService[*corev1.Pod]
	DeploymentService  *BaseResourceService[*appsv1.Deployment]
	DaemonSetService   *BaseResourceService[*appsv1.DaemonSet]
	IngressService     *BaseResourceService[*networkingv1.Ingress]
	ConfigMapService   *BaseResourceService[*corev1.ConfigMap]
	SecretService      *BaseResourceService[*corev1.Secret]
	PVCService         *BaseResourceService[*corev1.PersistentVolumeClaim]
	PVService          *BaseResourceService[*corev1.PersistentVolume]
	StatefulSetService *BaseResourceService[*appsv1.StatefulSet]
	NamespaceService   *BaseResourceService[*corev1.Namespace]

	// 非 K8s 服务
	InstallerService *InstallerService
	ProxyService     *ProxyService
}
