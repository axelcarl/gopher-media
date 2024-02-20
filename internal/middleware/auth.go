package middleware

import (
	"net/http"
	"strconv"

	"github.com/axelcarl/gopher-media/internal/cache"
	"github.com/axelcarl/gopher-media/internal/model"
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

func OwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The post id should be an unsigned integer.",
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		post, err := model.GetPost(uint(postID))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No post of that ID was found.",
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		value, exists := c.Get("userID")
		
		stringUserID, ok := value.(string)

		if !exists || !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
			})
			c.Abort()
			return
		}

    userID, err := strconv.ParseUint(stringUserID, 10, 64)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if uint(userID) != post.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "The post does not belong to this user.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
