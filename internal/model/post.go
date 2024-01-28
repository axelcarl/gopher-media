package model

import (
	"github.com/axelcarl/gopher-media/internal/database"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

func CreatePost(post *Post) error {
	result := database.DB.Create(&post)
	return result.Error
}

func GetPost(id uint) (*Post, error) {
	var post Post
	result := database.DB.First(&post, id)
	return &post, result.Error
}
