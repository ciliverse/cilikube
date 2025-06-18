package handlers

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse 定义了API错误响应的标准格式
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// respondError 返回错误响应
func respondError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// respondSuccess 返回成功响应
func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, SuccessResponse{
		Code:    code,
		Data:    data,
		Message: "success",
	})
}
