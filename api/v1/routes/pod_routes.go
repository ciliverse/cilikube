package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPodRoutes(router *gin.RouterGroup, handler *handlers.PodHandler) {
	// 基础资源操作
	podGroup := router.Group("/namespaces/:namespace/pods")
	{
		podGroup.GET("", handler.ListPods)
		podGroup.POST("", handler.CreatePod)
		podGroup.GET("/:name", handler.GetPod)
		podGroup.PUT("/:name", handler.UpdatePod)
		podGroup.DELETE("/:name", handler.DeletePod)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/pods")
	{
		watchGroup.GET("", handler.WatchPods)
	}
}
