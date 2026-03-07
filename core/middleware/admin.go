package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly returns a Gin middleware that restricts access to admin users.
// It reads the "role" value from the Gin context and aborts with 403 if
// the role is not "admin". Must be used after a middleware that sets
// c.Set("role", ...). Note: the built-in AuthMiddleware only sets "user_id";
// you must add your own middleware to load and set the user's role.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "admin access required",
		})
	}
}
