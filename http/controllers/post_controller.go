package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostController handles CRUD operations for posts.
type PostController struct{}

func (ctrl *PostController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController index"})
}

func (ctrl *PostController) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController create form"})
}

func (ctrl *PostController) Store(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "PostController store"})
}

func (ctrl *PostController) Show(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController show", "id": id})
}

func (ctrl *PostController) Edit(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController edit form", "id": id})
}

func (ctrl *PostController) Update(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "PostController update", "id": id})
}

func (ctrl *PostController) Destroy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PostController destroy"})
}
