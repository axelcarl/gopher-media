package handler

import (
	"fmt"
	"net/http"
	"strconv"
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
			"id":      userID,
		})
		return
	})

	router.GET("/authenticate", authFunc(rdb), func(c *gin.Context) {
		shouldReturnUser := c.Query("full")

		if shouldReturnUser != "true" {
			return
		}

		value, exists := c.Get("userID")

		stringUserID, ok := value.(string)

		if !exists || !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		userID, err := strconv.ParseUint(stringUserID, 10, 64)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		user, err := model.GetUser(uint(userID))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			return
		}

		c.JSON(http.StatusOK, user)
		return
	})
}
