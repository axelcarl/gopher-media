package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var users []string = []string{"Adam", "Tom", "Sara"}

func UserRoutes(router *gin.RouterGroup) {
	// Get endpoint /user/:id.
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer above 0.",
			})
			return
		}

		if id > len(users) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found.",
			})
			return
		}

		user := users[id-1]

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		return
	})

	// Post endpoint /user.
	router.POST("/", func(c *gin.Context) {
		users = append(users, "test")
	})

	// Put endpoint /user/:id.
	router.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		name := c.Query("name")

		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer above 0.",
			})
			return
		}

		if id > len(users) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found.",
			})
			return
		}

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "No name provided.",
			})
			return
		}

		users[id-1] = name

		c.JSON(http.StatusOK, gin.H{
			"user": users[id-1],
		})
		return
	})

	// Delete endpoint /user/:id.
	router.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer above 0.",
			})
			return
		}

		if id > len(users) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found.",
			})
			return
		}

		users[id-1] = "DELETED."

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted.",
		})
		return
	})
}
