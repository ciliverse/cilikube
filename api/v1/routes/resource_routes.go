package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
)

// RegisterResourceRoutes 注册资源路由
func RegisterResourceRoutes[T runtime.Object](router *gin.RouterGroup, handler *handlers.ResourceHandler[T], k8sManager *k8s.ClusterManager, resourceType string) {
	group := router.Group("/" + resourceType)
	{
		group.GET("", handler.List)
		group.POST("", handler.Create)
		group.GET("/:name", handler.Get)
		group.PUT("/:name", handler.Update)
		group.DELETE("/:name", handler.Delete)
		group.GET("/watch", handler.Watch)
	}
}
