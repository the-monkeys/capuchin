package main

import (
	"capuchin/internal/auth"
	"capuchin/internal/database"
	"capuchin/internal/todo"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present so os.Getenv reads values during Init
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}
	//Initialize database and auth
	auth.Init()
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

	// Public Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/signup", auth.Signup)
	r.POST("/login", auth.Login)

	// Protected Routes
	user := r.Group("/user")
	user.Use(auth.Middleware())
	{
		user.GET("/todos", todo.GetTodos)
		user.POST("/todos", todo.AddTodo)
		user.PATCH("/todos/:id", todo.ToggleTodo)
		user.PATCH("/todos/:id/edit", todo.EditTodo)
		user.DELETE("/todos/:id", todo.DeleteTodo)
	}

	r.Run(":8080")
}
