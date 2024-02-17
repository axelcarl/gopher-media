package middleware

import (
	"net/http"

	"github.com/axelcarl/gopher-media/internal/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, err := cache.GetUserIDBySession(rdb, sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid session.",
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}
