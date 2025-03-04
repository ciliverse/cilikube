package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterPVCRoutes 注册PVC相关路由
func RegisterPVCRoutes(router *gin.RouterGroup, handler *handlers.PVCHandler) {
	// 基础资源操作
	pvcGroup := router.Group("/namespaces/:namespace/pvcs")
	{
		pvcGroup.GET("", handler.ListPVCs)
		pvcGroup.POST("", handler.CreatePVC)
		pvcGroup.GET("/:name", handler.GetPVC)
		pvcGroup.PUT("/:name", handler.UpdatePVC)
		pvcGroup.DELETE("/:name", handler.DeletePVC)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/pvcs")
	{
		watchGroup.GET("", handler.WatchPVCs)
	}
}
