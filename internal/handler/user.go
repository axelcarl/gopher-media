package handler

import (
	"net/http"
	"strconv"

	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.RouterGroup) {
	// Get endpoint /user/:id.
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an unsigned integer.",
				"error":   err.Error(),
			})
			return
		}

		user, err := model.GetUser(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User does not exist.",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong",
					"error":   err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, user)
		return
	})

	// Post endpoint /user.
	router.POST("/", func(c *gin.Context) {
		var userCredentials model.UserRegistrationFields
		if err := c.ShouldBindJSON(&userCredentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request body is not a valid as a user.",
				"error":   err.Error(),
			})
			return
		}

		user := model.User{Name: userCredentials.Name, Password: userCredentials.Password}
		err := model.CreateUser(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User created.",
			"user_id": user.ID,
		})
		return
	})

	// Put endpoint /user/:id.
	router.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID must be an unsigned integer.",
				"error":   err.Error(),
			})
			return
		}

		var newFields model.UserUpdateFields
		if err := c.ShouldBindJSON(&newFields); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Body provided is not a valid user.",
				"error":   err.Error(),
			})
			return
		}

		user, err := model.GetUser(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User does not exist.",
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

		if err := model.UpdateUser(user, &newFields); err != nil {
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
			"message": "User was updated.",
		})
		return
	})

	// Delete endpoint /user/:id.
	router.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "ID must be an unsigned integer.",
				"error":   err.Error(),
			})
			return
		}

		user, err := model.GetUser(uint(id))

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User already deleted.",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong",
					"error":   err.Error(),
				})
			}
			return
		}

		if err = model.DeleteUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User was deleted.",
		})
		return
	})
}
