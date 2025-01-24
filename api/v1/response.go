package handlers

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Code: code,
		Data: data,
	})
}

func respondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{
		Code:    code,
		Message: message,
	})
}
