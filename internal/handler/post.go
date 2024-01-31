package handler

import (
	"net/http"
	"strconv"

	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.RouterGroup) {
	// Get endpoint /post/:id.
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID has to be of type uint.",
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post does not exist.",
			})
			return
		}

		c.JSON(http.StatusOK, post)
		return
	})

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
			"message": "Post created.",
			"post_id": post.ID,
			"title":   post.Title,
		})
		return
	})

	// Put endpoint /post:id
	router.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID has to be of type uint.",
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		var newFields model.PostUpdateFields

		err = c.ShouldBindJSON(&newFields)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Expected json object with mutable fields.",
			})
			return
		}

		err = model.UpdatePost(post, &newFields)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, post)
		return
	})

	// Get endpoint /post/:id.
	router.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID has to be of type uint.",
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post does not exist.",
			})
			return
		}

		err = model.DeletePost(post)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post deleted.",
		})
		return
	})
}
