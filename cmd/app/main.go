package main

import (
	"os"
	"time"

	"github.com/axelcarl/gopher-media/internal/cache"
	"github.com/axelcarl/gopher-media/internal/database"
	"github.com/axelcarl/gopher-media/internal/handler"
	"github.com/axelcarl/gopher-media/internal/middleware"
	"github.com/axelcarl/gopher-media/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// Setup CORS.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
				return origin == "https://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
}))

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
	handler.AuthRoutes(r.Group("/"), redisClient, middleware.AuthMiddleware)

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
