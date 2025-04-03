package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterPVRoutes 注册PV相关路由
func RegisterPVRoutes(router *gin.RouterGroup, handler *handlers.PVHandler) {
	// 基础资源操作
	pvGroup := router.Group("/pvs")
	{
		pvGroup.GET("", handler.ListPVs)
		pvGroup.POST("", handler.CreatePV)
		pvGroup.GET("/:name", handler.GetPV)
		pvGroup.PUT("/:name", handler.UpdatePV)
		pvGroup.DELETE("/:name", handler.DeletePV)
	}

	// Watch端点
	// watchGroup := router.Group("/watch/pvs")
	// {
	// 	watchGroup.GET("", handler.WatchPVs)
	// }
}
