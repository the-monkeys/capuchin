package main

import (
	"capuchin/internal/database"
	"capuchin/internal/handlers"
	"capuchin/internal/routes"
	"capuchin/internal/services"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	// Periodic cleanup prevents the revoked-token table from growing forever.
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			if err := database.CleanupTokens(); err != nil {
				log.Printf("Error cleaning up expired tokens: %v", err)
			}
		}
	}()

	// Handlers depend on interfaces so business logic can be swapped in tests.
	authService := services.NewAuthService()
	todoService := services.NewTodoService()

	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := gin.Default()

	// Allow cross-origin requests so a separately hosted frontend can call this API.
	// Restrict this in production to trusted origins.
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			// Short-circuit preflight checks to avoid running downstream handlers.
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Keep route wiring centralized so auth boundaries are easy to audit.
	routes.SetupRoutes(r, authHandler, todoHandler)

	r.Run(":8080")
}
