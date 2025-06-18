package initialization

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AppServices 结构体需要包含新增的 ClusterService。
type AppServices struct {
	ClusterService     *service.ClusterService
	ServiceService     *service.BaseResourceService[*corev1.Service]
	PodService         *service.BaseResourceService[*corev1.Pod]
	DeploymentService  *service.BaseResourceService[*appsv1.Deployment]
	DaemonSetService   *service.BaseResourceService[*appsv1.DaemonSet]
	IngressService     *service.BaseResourceService[*networkingv1.Ingress]
	ConfigMapService   *service.BaseResourceService[*corev1.ConfigMap]
	SecretService      *service.BaseResourceService[*corev1.Secret]
	PVCService         *service.BaseResourceService[*corev1.PersistentVolumeClaim]
	PVService          *service.BaseResourceService[*corev1.PersistentVolume]
	StatefulSetService *service.BaseResourceService[*appsv1.StatefulSet]
	NamespaceService   *service.BaseResourceService[*corev1.Namespace]
	InstallerService   service.InstallerService
}

// AppHandlers 结构体只保留必要的处理器
type AppHandlers struct {
	ClusterHandler   *handlers.ClusterHandler
	InstallerHandler *handlers.InstallerHandler
	AuthHandler      *handlers.AuthHandler
	ProxyHandler     *handlers.ProxyHandler
}

// InitializeServices 初始化所有服务
func InitializeServices(k8sManager *k8s.ClusterManager, cfg *configs.Config) *AppServices {
	log.Println("正在初始化服务层...")

	resourceFactory := service.NewResourceServiceFactory()
	resourceFactory.InitializeDefaultServices()

	appServices := &AppServices{
		ClusterService:   service.NewClusterService(k8sManager),
		InstallerService: service.NewInstallerService(cfg),
	}

	initializeResourceService(resourceFactory, "services", &appServices.ServiceService)
	initializeResourceService(resourceFactory, "pods", &appServices.PodService)
	initializeResourceService(resourceFactory, "deployments", &appServices.DeploymentService)
	initializeResourceService(resourceFactory, "daemonsets", &appServices.DaemonSetService)
	initializeResourceService(resourceFactory, "ingresses", &appServices.IngressService)
	initializeResourceService(resourceFactory, "configmaps", &appServices.ConfigMapService)
	initializeResourceService(resourceFactory, "secrets", &appServices.SecretService)
	initializeResourceService(resourceFactory, "persistentvolumeclaims", &appServices.PVCService)
	initializeResourceService(resourceFactory, "persistentvolumes", &appServices.PVService)
	initializeResourceService(resourceFactory, "statefulsets", &appServices.StatefulSetService)
	initializeResourceService(resourceFactory, "namespaces", &appServices.NamespaceService)

	return appServices
}

// initializeResourceService is a helper function to initialize a specific resource service.
func initializeResourceService[T runtime.Object](factory *service.ResourceServiceFactory, resourceName string, serviceField **service.BaseResourceService[T]) {
	if svc, ok := factory.GetService(resourceName).(*service.BaseResourceService[T]); ok {
		*serviceField = svc
	} else {
		log.Fatalf("Failed to initialize %s service: type assertion failed or service not found", resourceName)
	}
}

// InitializeHandlers 初始化所有处理器
func InitializeHandlers(router *gin.RouterGroup, services *AppServices, k8sManager *k8s.ClusterManager) {
	// 注册集群相关路由
	clusterHandler := handlers.NewClusterHandler(services.ClusterService)
	clusterGroup := router.Group("/clusters")
	{
		clusterGroup.GET("", clusterHandler.ListClusters)
		clusterGroup.POST("", clusterHandler.CreateCluster)
		clusterGroup.DELETE("/:name", clusterHandler.DeleteCluster)
		clusterGroup.POST("/active", clusterHandler.SetActiveCluster)
		clusterGroup.GET("/active", clusterHandler.GetActiveCluster)
	}

	// In InitializeHandlers, the services are now typed correctly in the AppServices struct.
	// So, we can directly use them.
	registerResourceHandler(router, services.ServiceService, k8sManager, "services")
	registerResourceHandler(router, services.PodService, k8sManager, "pods")
	registerResourceHandler(router, services.DeploymentService, k8sManager, "deployments")
	registerResourceHandler(router, services.DaemonSetService, k8sManager, "daemonsets")
	registerResourceHandler(router, services.IngressService, k8sManager, "ingresses")
	registerResourceHandler(router, services.ConfigMapService, k8sManager, "configmaps")
	registerResourceHandler(router, services.SecretService, k8sManager, "secrets")
	registerResourceHandler(router, services.PVCService, k8sManager, "persistentvolumeclaims")
	registerResourceHandler(router, services.PVService, k8sManager, "persistentvolumes")
	registerResourceHandler(router, services.StatefulSetService, k8sManager, "statefulsets")
	registerResourceHandler(router, services.NamespaceService, k8sManager, "namespaces")

	// 注册安装器路由
	installerHandler := handlers.NewInstallerHandler(services.InstallerService)
	installerGroup := router.Group("/installers")
	{
		installerGroup.GET("/stream", installerHandler.StreamMinikubeInstallation)
	}

	// 注册代理路由
	proxyHandler := handlers.NewProxyHandler(k8sManager)
	routes.KubernetesProxyRoutes(router, proxyHandler)
}

// registerResourceHandler is a helper function to register a resource handler and its routes.
func registerResourceHandler[T runtime.Object](router *gin.RouterGroup, resourceService service.ResourceService[T], k8sManager *k8s.ClusterManager, resourceName string) {
	if resourceService != nil {
		handler := handlers.NewResourceHandler(resourceService, k8sManager, resourceName)
		routes.RegisterResourceRoutes(router, handler, k8sManager, resourceName)
	}
}

// SetupRouter 设置路由
func SetupRouter(cfg *configs.Config, services *AppServices, k8sManager *k8s.ClusterManager, e *casbin.Enforcer) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// API 路由组
	api := router.Group("/api")
	v1 := api.Group("/v1")

	// 初始化处理器
	InitializeHandlers(v1, services, k8sManager)

	return router
}
