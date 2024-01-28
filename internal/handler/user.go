package handler

import (
	"net/http"
	"strconv"

	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	// Get endpoint /user/:id.
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer above 0.",
			})
			return
		}

		user, err := model.GetUser(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":  user.ID,
			"username": user.Name,
		})
		return
	})

	// Post endpoint /user.
	router.POST("/", func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := model.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "User created.",
			"user_id":  user.ID,
			"username": user.Name,
		})
		return
	})

	// Put endpoint /user/:id.
	router.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID must be an unsigned integer.",
			})
			return
		}

		var newFields model.UserUpdateFields
		if err := c.ShouldBindJSON(&newFields); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Body provided is not a valid user.",
			})
			return
		}

		oldUser, err := model.GetUser(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found.",
			})
			return
		}

		if err := model.UpdateUser(oldUser, &newFields); err != nil {
			if err.Error() == "No valid field provided." {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "No valid fields provided.",
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

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated",
		})
		return
	})

	// Delete endpoint /user/:id.
	router.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID must be a valid unsigned integer above 0.",
			})
			return
		}

		if err := model.DeleteUser(uint(id)); err != nil {
			if err.Error() == "User not found." {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "User not found.", 
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong.",
					"error": err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted.",
		})
		return
	})
}
