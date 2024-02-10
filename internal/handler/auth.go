package handler

import (
	"net/http"

	"github.com/axelcarl/gopher-media/internal/model"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	router.POST("/login", func(c *gin.Context) {
		var credentials model.UserLoginFields
		err := c.ShouldBindJSON(&credentials)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Expected username and password",
				"error":   err.Error(),
			})
			return
		}

		err = model.LoginUser(&credentials)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Passwords are not matching.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in!",
		})
		return
	})
}
