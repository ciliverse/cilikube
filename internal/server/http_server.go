package server

import (
	v1 "github.com/ciliverse/cilikube/api/v1"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	v1.RegisterRoutes(r)
	r.Run(":8080") // 启动服务器
}
