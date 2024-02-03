package model

import (
	"errors"

	"github.com/axelcarl/gopher-media/internal/database"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Text   string `json:"text" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type PostUpdateFields struct {
	Title *string `json:"title" binding:"required"`
	Text  *string `json:"text" binding:"required"`
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

func UpdatePost(post *Post, postUpdateFields *PostUpdateFields) error {
	if postUpdateFields.Text != nil && postUpdateFields.Title != nil {
		post.Title = *postUpdateFields.Title
		post.Text = *postUpdateFields.Text
		result := database.DB.Save(&post)
		return result.Error
	}
	return errors.New("Missing fields.")
}

func DeletePost(post *Post) error {
	result := database.DB.Delete(&post)
	return result.Error
}
