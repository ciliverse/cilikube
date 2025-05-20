// api/v1/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 定义了成功响应的通用结构
type SuccessResponse struct {
	Code    int         `json:"code"`           // HTTP 状态码或业务成功码
	Message string      `json:"message"`        // 成功信息
	Data    interface{} `json:"data,omitempty"` // 实际数据, omitempty 表示如果为空则不显示此字段
}

// ErrorResponse 定义了错误响应的通用结构
type ErrorResponse struct {
	Code    int         `json:"code"`              // HTTP 状态码或业务错误码
	Message string      `json:"message"`           // 错误描述
	Details interface{} `json:"details,omitempty"` // 错误的详细信息 (可选)
}

// SendSuccess 是一个帮助函数，用于发送成功的 HTTP 响应
// 它使用 HTTP 200 OK 状态码。
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK, // 你也可以定义自己的业务成功码
		Message: message,
		Data:    data,
	})
}

// SendError 是一个帮助函数，用于发送错误的 HTTP 响应
// statusCode 是实际的 HTTP 状态码 (例如 http.StatusBadRequest, http.StatusNotFound)。
func SendError(c *gin.Context, statusCode int, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode, // 通常错误码与 HTTP 状态码一致
		Message: message,
		Details: details,
	})
}

// SendRawSuccess 发送成功的 HTTP 响应，允许自定义成功码
func SendRawSuccess(c *gin.Context, httpStatusCode int, successCode int, message string, data interface{}) {
	c.JSON(httpStatusCode, SuccessResponse{
		Code:    successCode,
		Message: message,
		Data:    data,
	})
}

// SendRawError 发送错误的 HTTP 响应，允许自定义错误码
func SendRawError(c *gin.Context, httpStatusCode int, errorCode int, message string, details interface{}) {
	c.JSON(httpStatusCode, ErrorResponse{
		Code:    errorCode,
		Message: message,
		Details: details,
	})
}
