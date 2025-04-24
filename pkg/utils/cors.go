package utils

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 处理跨域请求，支持预检请求和 SSE
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的 Origin
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{
			"http://192.168.1.100:8888", // 明确指定前端地址
			"http://localhost:8888",     // 本地开发环境（可选）
			"http://192.168.1.1",
		}

		// 检查请求的 Origin 是否在允许列表中
		allowedOrigin := ""
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = origin
				break
			}
		}

		// 如果 Origin 不匹配，fallback 到默认行为（不设置 CORS 头）
		if allowedOrigin == "" {
			log.Printf("CORS: Origin %s not allowed, skipping CORS headers", origin)
			c.Next()
			return
		}

		// 设置跨域相关响应头
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true") // 支持凭据
		c.Header("Access-Control-Max-Age", "86400")          // 预检请求缓存时间

		// 根据请求路径设置 Content-Type
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/system/install-minikube") {
			c.Header("Content-Type", "text/event-stream; charset=utf-8")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
		} else {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}

		// 添加日志输出
		log.Printf("CORS headers set for %s %s (Origin: %s): %v",
			c.Request.Method, c.Request.URL.Path, origin, c.Writer.Header())

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
