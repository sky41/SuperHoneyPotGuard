package middleware

import (
	"net/http"
	"strings"
	"superhoneypotguard/models"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ErrorResponse(c, http.StatusUnauthorized, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) *models.Claims {
	userId, _ := c.Get("userId")
	username, _ := c.Get("username")
	roles, _ := c.Get("roles")
	permissions, _ := c.Get("permissions")

	return &models.Claims{
		UserID:      userId.(int),
		Username:    username.(string),
		Roles:       roles.([]string),
		Permissions: permissions.([]string),
	}
}
