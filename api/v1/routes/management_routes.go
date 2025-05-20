// api/v1/routes/management_routes.go
package routes

import (
	"log"

	"github.com/ciliverse/cilikube/api/v1/handlers"
	// "github.com/ciliverse/cilikube/pkg/auth" // 认证中间件

	"github.com/gin-gonic/gin"
)

// RegisterManagementRoutes 注册与平台管理功能相关的 API 路由。
// managementHandler 参数是已初始化的 ManagementHandler 实例。
func RegisterManagementRoutes(routerGroup *gin.RouterGroup, managementHandler *handlers.ManagementHandler) {
	if managementHandler == nil {
		log.Println("警告: ManagementHandler 为 nil，无法注册 Management 相关路由。")
		return
	}

	// 创建管理接口的路由组，例如 /api/v1/management
	managementSpecificGroup := routerGroup.Group("/management")
	// managementSpecificGroup.Use(auth.JWTMiddleware(), auth.AdminOnlyMiddleware()) // 示例：更严格的访问控制

	{
		// GET /api/v1/management/clusters - 列出所有可用的集群名称
		managementSpecificGroup.GET("/clusters", managementHandler.ListAvailableClusters)
		log.Printf("路由注册: GET %s/management/clusters", routerGroup.BasePath())

		// 未来可以添加其他管理接口，例如：
		// POST /api/v1/management/clusters - 用于动态添加新集群（如果实现了此功能）
		// DELETE /api/v1/management/clusters/:clusterName - 用于移除集群（如果实现了此功能）
	}
	log.Println("Management 相关路由注册完成。")
}
