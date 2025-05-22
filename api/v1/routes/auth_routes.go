package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	// 认证路由组
	authGroup := router.Group("/api/v1/auth")
	{
		// 公开路由（不需要认证）
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)

		// 需要认证的路由
		authenticated := authGroup.Group("")
		authenticated.Use(auth.JWTAuthMiddleware())
		{
			authenticated.GET("/profile", authHandler.GetProfile)
			authenticated.PUT("/profile", authHandler.UpdateProfile)
			authenticated.POST("/change-password", authHandler.ChangePassword)
			authenticated.POST("/logout", authHandler.Logout)
		}

		// 管理员专用路由
		admin := authGroup.Group("")
		admin.Use(auth.JWTAuthMiddleware(), auth.AdminRequiredMiddleware())
		{
			admin.GET("/users", authHandler.GetUserList)
			admin.PUT("/users/:id/status", authHandler.UpdateUserStatus)
			admin.DELETE("/users/:id", authHandler.DeleteUser)
		}
	}
}
