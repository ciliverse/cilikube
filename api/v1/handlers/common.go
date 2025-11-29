package handlers

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse defines the standard format for API error responses
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse defines the structure for successful responses
type SuccessResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// respondError returns an error response
func respondError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// respondSuccess returns a successful response
func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, SuccessResponse{
		Code:    code,
		Data:    data,
		Message: "success",
	})
}
