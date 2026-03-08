package routes

import (
	"capuchin/internal/handlers"
	"capuchin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// Global Error Handler
	router.Use(middleware.ErrorHandler())

	// Public Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	// Protected Routes
	protected := router.Group("/api/user")
	protected.Use(middleware.AuthRequired()) // Using the new middleware package
	{
		protected.GET("/todo", handlers.GetTodos)            // Get all todos for the authenticated user
		protected.POST("/todo", handlers.AddTodo)            // Add a new todo for the authenticated user
		protected.PATCH("/todo/:id", handlers.ToggleTodo)    // Toggle the completion status of a specific todo
		protected.PATCH("/todo/:id/edit", handlers.EditTodo) // Edit the content of a specific todo
		protected.DELETE("/todo/:id", handlers.DeleteTodo)   // Delete a specific todo
	}
}
