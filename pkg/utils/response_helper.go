package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard API response helper
func ApiSuccess(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    data,
		"message": message,
	})
}

func ApiError(c *gin.Context, statusCode int, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}
	log.Printf("API Error: Status %d, Message: %s, Details: %s, Path: %s", statusCode, message, detailStr, c.Request.URL.Path)
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"data":    nil,
		"message": message,
		"details": detailStr,
	})
}
