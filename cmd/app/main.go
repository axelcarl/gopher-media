package main

import (
	"os"

	"github.com/axelcarl/gopher-media/internal/database"
	"github.com/axelcarl/gopher-media/internal/handler"
	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// load dotenv variables.
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables.")
	}

	// Create a database connection.
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to establish a database connection.")
	}

	// Initialize the database.
	database.InitDatabase(db)

	// Run migrations.
	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Post{})

	// Setup routes.
	handler.AuthRoutes(r.Group("/"))
	handler.UserRoutes(r.Group("/user"))
	handler.PostRoutes(r.Group("/post"))

	// Run application.
	r.Run()
}
