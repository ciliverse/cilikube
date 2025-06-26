package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 处理跨域请求，支持预检请求
// 建议: 将 allowedOrigins 作为参数传入或从配置加载
func Cors(allowedOrigins []string) gin.HandlerFunc {
	// 如果 allowedOrigins 为空，提供一个默认值或打印警告
	if len(allowedOrigins) == 0 {
		log.Println("CORS Warning: No allowed origins configured!")
		// 可以选择完全禁用 CORS 或允许所有 (如果 AllowCredentials 为 false)
		// 这里选择禁用 CORS，只继续处理
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// 获取请求的 Origin
		origin := c.Request.Header.Get("Origin")

		// 如果没有 Origin 头（例如，非浏览器请求或同源请求），则无需处理 CORS
		if origin == "" {
			c.Next()
			return
		}

		// 检查请求的 Origin 是否在允许列表中
		allowedOrigin := ""
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = origin
				break
			}
			// 可选: 处理通配符子域名等更复杂的匹配逻辑
		}

		// 如果 Origin 匹配成功
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			// 重要：因为 Allow-Origin 不是 "*"，所以需要 Vary 头
			c.Header("Vary", "Origin")

			// 设置其他 CORS 头
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")            // 更全面的方法列表
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, Accept") // 添加常用头
			c.Header("Access-Control-Allow-Credentials", "true")                                          // 支持凭据
			c.Header("Access-Control-Max-Age", "86400")                                                   // 预检请求缓存时间

			// 正确处理 OPTIONS 预检请求：如果 Origin 允许，且方法是 OPTIONS，则中止并返回 204
			if c.Request.Method == "OPTIONS" {
				log.Printf("CORS: Preflight request for %s from %s allowed.", c.Request.URL.Path, origin)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			// 对于非 OPTIONS 请求，继续处理
			log.Printf("CORS: Allowed non-preflight request for %s from %s.", c.Request.URL.Path, origin)
			c.Next()

		} else {
			// 如果 Origin 不匹配
			log.Printf("CORS: Origin '%s' not allowed for %s.", origin, c.Request.URL.Path)

			// 对于 OPTIONS 预检请求，如果 Origin 不允许，也应该中止，但可以不设置 CORS 头
			// 浏览器会因为缺少必要的 Allow-Origin 头而拒绝该请求
			if c.Request.Method == "OPTIONS" {
				// 可以选择返回 403 Forbidden 或 просто Abort
				// c.AbortWithStatus(http.StatusForbidden) // 更明确的拒绝
				c.Abort() // 或者仅仅中止处理链
				return
			}

			// 对于非 OPTIONS 请求，Origin 不允许
			// 浏览器会发出请求，但会阻止前端 JS 读取响应。
			// 在这里可以选择：
			// 1. 调用 c.Next()：让请求继续，但浏览器会报错（当前你的代码逻辑）
			// 2. 调用 c.AbortWithStatus(http.StatusForbidden)：直接拒绝请求（更严格）
			// 选择 c.Next() 意味着后端可能执行了操作，但前端无法收到结果。
			// 选择 Abort 更安全，阻止未经授权的源触发操作。
			// 这里我们选择更安全的 Abort:
			log.Printf("CORS: Aborting non-preflight request from disallowed origin '%s' for %s.", origin, c.Request.URL.Path)
			c.AbortWithStatus(http.StatusForbidden) // 直接拒绝
			// 或者，如果你想保持原来的行为（允许后端处理但浏览器阻止）:
			// c.Next()

			return // 确保中止后返回
		}
	}
}

// 在 main.go 中使用修改后的自定义中间件
/*
func main() {
    router := gin.Default()

    // 从配置或环境变量加载允许的源
    origins := []string{
        "http://192.168.1.100:8888",
        "http://localhost:8888",
        "http://192.168.1.1",
    }
    router.Use(utils.Cors(origins)) // 使用你的自定义 CORS 中间件

    // ... 注册路由 ...

    router.Run(":8080")
}
*/
