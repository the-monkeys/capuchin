package main

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"capuchin/internal/handlers"
	"capuchin/internal/middleware"
	"capuchin/internal/routes"
	"capuchin/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect(config.Config)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := database.CleanupTokens(); err != nil {
				log.Printf("token cleanup: error: %v", err)
			}
		}
	}()

	authService := services.NewAuthService()
	todoService := services.NewTodoService()

	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := gin.Default()

	// Restrict Access-Control-Allow-Origin to trusted origins in production.
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Use(middleware.DBHealthCheck())

	routes.SetupRoutes(r, authHandler, todoHandler)

	r.Run(":8080")
}
