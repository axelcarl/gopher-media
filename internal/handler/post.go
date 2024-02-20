package handler

import (
	"net/http"
	"strconv"

	"github.com/axelcarl/gopher-media/internal/middleware"
	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRoutes(router *gin.RouterGroup) {
	// Get endpoint /post/:id.
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an unsigned integer.",
				"error":   err.Error(),
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Post does not exist.",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong.",
					"error":   err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, post)
		return
	})

	// Post endpoint /post.
	router.POST("/", func(c *gin.Context) {
		var post model.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request body is not a valid post.",
				"error":   err.Error(),
			})
			return
		}

		// Look if user exists.
		_, err := model.GetUser(post.UserID)

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Creator of post (user) does not exist.",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong.",
					"error":   err.Error(),
				})
			}
			return
		}

		err = model.CreatePost(&post)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post created.",
			"post_id": post.ID,
		})
		return
	})

	// Put endpoint /post:id
	router.PUT("/:id", middleware.OwnerMiddleware(), func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID has to be of type uint.",
				"error":   err.Error(),
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Post does not exist",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong.",
					"error":   err.Error(),
				})
			}
			return
		}

		var newFields model.PostUpdateFields

		err = c.ShouldBindJSON(&newFields)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Body provided is not a valid post.",
				"error":   err.Error(),
			})
			return
		}

		err = model.UpdatePost(post, &newFields)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post was updated.",
		})
		return
	})

	// Get endpoint /post/:id.
	router.DELETE("/:id", middleware.OwnerMiddleware(), func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID must be an unsigned integer.",
				"error":   err.Error(),
			})
			return
		}

		post, err := model.GetPost(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Post already deleted.",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong.",
					"error":   err.Error(),
				})
			}
			return
		}

		err = model.DeletePost(post)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post was deleted.",
		})
		return
	})
}
