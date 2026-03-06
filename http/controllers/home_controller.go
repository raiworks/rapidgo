package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home renders the home page using the home.html template.
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Welcome to RGo",
	})
}
