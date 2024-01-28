package handler

import (
	"net/http"

	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.RouterGroup) {
	// Post endpoint /post.
	router.POST("/", func(c *gin.Context) {
		var post model.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Look if user exists.
		_, err := model.GetUser(post.UserID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "User ID provided is not valid.",
			})
			return
		}

		err = model.CreatePost(&post)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "User created.",
			"post_id":  post.ID,
			"title": post.Title,
		})
		return
	})

}
