package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// 设置为发布模式，关闭调试信息
	gin.SetMode(gin.ReleaseMode)

	// 创建路由引擎
	r := gin.Default()

	// 设置受信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin!",
		})
	})

	// 启动服务器
	r.Run(":8080")
}
