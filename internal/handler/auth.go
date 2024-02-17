package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/axelcarl/gopher-media/internal/cache"
	"github.com/axelcarl/gopher-media/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const ONE_DAY int = 86400

func AuthRoutes(router *gin.RouterGroup, rdb *redis.Client, authFunc func(rdb *redis.Client) gin.HandlerFunc) {
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

		userID, err := model.LoginUser(&credentials)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Passwords are not matching.",
				"error":   err.Error(),
			})
			return
		}

		sessionID, _ := cache.CreateSession(rdb, fmt.Sprint(userID), time.Hour)

		c.SetCookie("session_id", sessionID, ONE_DAY, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in.",
		})
		return
	})

	router.GET("/authenticate", authFunc(rdb))
}
