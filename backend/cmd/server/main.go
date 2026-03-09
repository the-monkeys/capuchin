package main

import (
	"capuchin/internal/database"
	"capuchin/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	//Initialize database
	database.Connect()
	database.InitSchema()

	//Initialize Gin router
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Inject all predefined routes
	routes.SetupRoutes(r)

	r.Run(":8080")
}
