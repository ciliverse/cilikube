package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// [核心修改] 函数名改为 RegisterAuthRoutes 并接收 *gin.RouterGroup
// 这样它就和其他 Register... 函数保持一致了
func RegisterAuthRoutes(authGroup *gin.RouterGroup) {
	authHandler := handlers.NewAuthHandler()

	// 路由直接注册在传入的 authGroup 上，不再自己创建

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
	admin := authGroup.Group("/admin") // 将管理路由分组到 /admin 下更清晰
	admin.Use(auth.JWTAuthMiddleware(), auth.AdminRequiredMiddleware())
	{
		admin.GET("/users", authHandler.GetUserList)
		admin.PUT("/users/:id/status", authHandler.UpdateUserStatus)
		admin.DELETE("/users/:id", authHandler.DeleteUser)
	}
}
