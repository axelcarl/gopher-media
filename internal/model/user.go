package model

import (
	"errors"

	"github.com/axelcarl/gopher-media/internal/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Posts    []Post `json:"posts"`
	Password string `json:"-"`
}

type UserRegistrationFields struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateFields struct {
	Name *string `json:"name" binding:"required"`
}

type UserLoginFields struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordBytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(user *User) error {
	hashedPassword, err := HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := database.DB.Create(&user)
	return result.Error
}

func LoginUser(credentials *UserLoginFields) (uint, error) {
	var user User
	result := database.DB.Where("name = ?", credentials.Name).First(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	match := CheckPasswordHash(user.Password, credentials.Password)

	if match == false {
		return 0, errors.New("Passwords not matching.")
	}

	return user.ID, nil
}

func GetUser(id uint) (*User, error) {
	var user User
	result := database.DB.Preload("Posts").First(&user, id)
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

func DeleteUser(user *User) error {
	result := database.DB.Delete(&user)
	return result.Error
}
