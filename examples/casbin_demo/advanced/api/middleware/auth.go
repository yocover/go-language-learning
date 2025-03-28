package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 创建一个权限验证中间件
func AuthMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头或 JWT 获取用户信息
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未认证",
			})
			c.Abort()
			return
		}

		// 获取域信息（可以从请求头或其他地方获取）
		domain := c.GetHeader("X-Domain")
		if domain == "" {
			domain = "platform" // 默认域
		}

		// 获取请求方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path

		// 检查权限
		allowed, err := enforcer.Enforce(userID, domain, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "权限检查失败",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "没有权限",
			})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", userID)
		c.Set("domain", domain)

		c.Next()
	}
}
