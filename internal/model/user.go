package model

import (
	"errors"

	"github.com/axelcarl/gopher-media/internal/database"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Posts []Post `json:"post" binding:"required"`
}

type UserUpdateFields struct {
	Name *string `json:"name" binding:"required"`
}

func CreateUser(user *User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func GetUser(id uint) (*User, error) {
	var user User
	result := database.DB.First(&user, id)
	return &user, result.Error
}

func UpdateUser(user *User, newFields *UserUpdateFields) error {
	if newFields.Name != nil {
		user.Name = *newFields.Name
		result := database.DB.Model(&user).Update("name", user.Name)
		return result.Error
	}

	return errors.New("No valid field provided.")
}

func DeleteUser(id uint) error {
	user, err := GetUser(id)
	if err != nil {
		return errors.New("User not found.")
	}

	result := database.DB.Delete(&user)
	return result.Error
}
