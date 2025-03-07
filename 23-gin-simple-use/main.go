package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的gin路由引擎
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin!",
		})
	})

	// 带参数的路由
	r.GET("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name,
		})
	})

	// 带查询参数的路由
	r.GET("/search", func(ctx *gin.Context) {
		q := ctx.Query("q")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Searching for: " + q,
		})
	})

	// POST 请求的路由
	r.POST("/submit", func(ctx *gin.Context) {
		var json struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Email string `json:"email"`
		}
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User submitted successfully",
			"data":    json,
		})
	})

	r.Run(":8080")
}
