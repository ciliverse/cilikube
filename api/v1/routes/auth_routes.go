package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, handler *handlers.AuthHandler) {
	// 注册路由
	auth := router.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/create", handler.CreateUser)
	}
}
