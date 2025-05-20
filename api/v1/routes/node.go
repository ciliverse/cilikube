// api/v1/routes/node_routes.go
package routes

import (
	"log"

	"github.com/ciliverse/cilikube/api/v1/handlers"
	// "github.com/ciliverse/cilikube/pkg/auth" // 如果有认证中间件

	"github.com/gin-gonic/gin"
)

// RegisterNodeRoutes 注册与节点相关的 API 路由。
// 这些路由现在是集群特定的，路径中包含 :clusterName。
// nodeHandler 参数是已初始化的 NodeHandler 实例。
func RegisterNodeRoutes(routerGroup *gin.RouterGroup, nodeHandler *handlers.NodeHandler) {
	if nodeHandler == nil {
		log.Println("警告: NodeHandler 为 nil，无法注册 Node 相关路由。")
		return
	}

	// 创建一个专门用于处理特定集群下节点资源的路由组
	// 例如: /api/v1/clusters/{clusterName}/nodes
	// 这里的 routerGroup 应该是 /api/v1
	// 所以实际路径会是 /api/v1/clusters/:clusterName/...
	nodesSpecificClusterGroup := routerGroup.Group("/clusters/:clusterName")
	// nodesSpecificClusterGroup.Use(auth.JWTMiddleware()) // 示例：应用认证中间件

	{
		// GET /api/v1/clusters/:clusterName/nodes
		nodesSpecificClusterGroup.GET("/nodes", nodeHandler.ListNodes)
		log.Printf("路由注册: GET %s/clusters/:clusterName/nodes", routerGroup.BasePath())

		// GET /api/v1/clusters/:clusterName/nodes/:nodeName
		nodesSpecificClusterGroup.GET("/nodes/:nodeName", nodeHandler.GetNodeDetails)
		log.Printf("路由注册: GET %s/clusters/:clusterName/nodes/:nodeName", routerGroup.BasePath())

		// 未来可以添加其他与特定集群内节点相关的路由
	}
	log.Println("Node 相关路由注册完成。")
}
