package middleware

import (
	"net/http"
	"strings"

	"github.com/RAiWorks/RGo/core/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT Bearer tokens on protected routes.
// On success, sets "user_id" in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
