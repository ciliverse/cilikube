package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由器
	router := gin.New()
	router.GET("/health", HealthCheck)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "timestamp")
	assert.Contains(t, w.Body.String(), "uptime")
}

func TestReadinessCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/ready", ReadinessCheck)

	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
	assert.Contains(t, w.Body.String(), "checks")
}

func TestLivenessCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/live", LivenessCheck)

	req, _ := http.NewRequest("GET", "/live", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
}
