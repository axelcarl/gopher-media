package main

import (
	"github.com/axelcarl/gopher-media/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	handler.UserRoutes(r.Group("/user"))

	r.Run()
}
