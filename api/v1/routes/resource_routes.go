package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterResourceRoutes 通用路由注册
func RegisterResourceRoutes[T any](router *gin.RouterGroup, handler *handlers.ResourceHandler[T], resource string) {
	group := router.Group("/" + resource)
	group.GET("", handler.List)
	group.POST("", handler.Create)
	group.GET("/:name", handler.Get)
	group.PUT("/:name", handler.Update)
	group.DELETE("/:name", handler.Delete)
}
