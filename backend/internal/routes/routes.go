package routes

import (
	"capuchin/internal/handlers"
	"capuchin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler) {

	// Global Error Handler
	router.Use(middleware.ErrorHandler())

	// Public Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.POST("/signup", authHandler.Signup)
	router.POST("/login", authHandler.Login)

	// Protected Routes
	protected := router.Group("/api/user")
	protected.Use(middleware.AuthRequired()) // Using the new middleware package
	{
		protected.POST("/logout", authHandler.Logout)           // Logout the authenticated user
		protected.GET("/todo", todoHandler.GetTodos)            // Get all todos for the authenticated user
		protected.POST("/todo", todoHandler.AddTodo)            // Add a new todo for the authenticated user
		protected.PATCH("/todo/:id", todoHandler.ToggleTodo)    // Toggle the completion status of a specific todo
		protected.PATCH("/todo/:id/edit", todoHandler.EditTodo) // Edit the content of a specific todo
		protected.DELETE("/todo/:id", todoHandler.DeleteTodo)   // Delete a specific todo
	}
}
