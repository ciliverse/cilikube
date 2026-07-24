package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
)

// GetShowcaseInfo is a public endpoint. Credentials are returned only when
// CILIKUBE_SHOWCASE=1 (fail-closed otherwise).
func GetShowcaseInfo(c *gin.Context) {
	info := k8s.PublicShowcaseInfo()
	if !info.Showcase {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    gin.H{"showcase": false},
			"message": "ok",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    info,
		"message": "showcase demo info",
	})
}
