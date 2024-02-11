package main

import (
	"os"

	"github.com/axelcarl/gopher-media/internal/cache"
	"github.com/axelcarl/gopher-media/internal/database"
	"github.com/axelcarl/gopher-media/internal/handler"
	"github.com/axelcarl/gopher-media/internal/middleware"
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

	// Initialize redis.
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDatabase := os.Getenv("REDIS_DATABASE")
	redisClient, err := cache.InitRedis(redisUrl, redisPassword, redisDatabase)

	if err != nil {
		panic("Failed to initialize redis.")
	}

	// Initialize the database.
	database.InitDatabase(db)

	// Run migrations.
	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Post{})

	// Setup routing.

	// Auth routes.
	handler.AuthRoutes(r.Group("/"), redisClient)

	// User routes.
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(redisClient))
	handler.UserRoutes(userRoutes)

	// Post routes.
	postRoutes := r.Group("/post")
	postRoutes.Use(middleware.AuthMiddleware(redisClient))
	handler.PostRoutes(postRoutes)

	// Run application.
	r.Run()
}
