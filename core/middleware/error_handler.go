package middleware

import (
	"github.com/raiworks/rapidgo/v2/core/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler returns middleware that processes errors added to the Gin context.
// If an error is an *errors.AppError, it uses the status code and ErrorResponse().
// Other errors are treated as 500 Internal Server Error.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		lastErr := c.Errors.Last().Err

		if appErr, ok := lastErr.(*errors.AppError); ok {
			c.JSON(appErr.Status, appErr.ErrorResponse())
			return
		}

		wrapped := errors.Internal(lastErr)
		c.JSON(wrapped.Status, wrapped.ErrorResponse())
	}
}
