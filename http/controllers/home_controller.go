package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home handles the root route and returns a welcome message.
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to RGo",
	})
}
