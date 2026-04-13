package routes

import (
	"capuchin/internal/handlers"
	"capuchin/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler) {

	router.Use(middleware.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/signup", authHandler.Signup)
	router.POST("/login", authHandler.Login)

	// Group authenticated routes so auth middleware is applied consistently.
	protected := router.Group("/api/user")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/todo", todoHandler.GetTodos)
		protected.POST("/todo", todoHandler.AddTodo)
		protected.PATCH("/todo/:id", todoHandler.UpdateTodo)
		protected.DELETE("/todo/:id", todoHandler.DeleteTodo)
	}
}
