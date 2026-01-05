package main

import (
	"capuchin/internal/models"
	"capuchin/internal/store"
	"net/http"

	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// path to db
const dbPath = "../../db/db.json"

// create a list of todos imported from models/todo.go
var todos []models.Todo

func main() {
	// Load data from disk on startup
	var err error
	todos, err = store.Load(dbPath)
	if err != nil {
		//if loading  fails on first run (file not found), just start empty
		todos = []models.Todo{}
	}

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	//health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// route to return  list of todos
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(200, todos)
	})

	// route to add a new item
	r.POST("/todos", func(c *gin.Context) {
		var newTodo models.Todo

		// Bind the incoming JSON to our struct
		if err := c.ShouldBindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newTodo.ID = uuid.New().String()

		// Add to the list
		todos = append(todos, newTodo)
		store.Save(dbPath, todos) // Save
		// Respond with the created item
		c.JSON(http.StatusOK, newTodo)
	})

	// PATCH /todos/:id - Toggle "completed" status
	r.PATCH("/todos/:id", func(c *gin.Context) {
		id := c.Param("id") // Get the ID from the URL
		// Iterate through the list to find the item
		for i, t := range todos {
			if t.ID == id {
				// Flip the status
				todos[i].Completed = !todos[i].Completed

				// Respond with the updated item
				c.JSON(200, todos[i])
				store.Save(dbPath, todos) // <--- Use store.Save
				return
			}
		}

		c.JSON(404, gin.H{"message": "Todo not found"})
	})

	// DELETE /todos/:id - Delete an item
	r.DELETE("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, t := range todos {
			if t.ID == id {
				// Delete: Append everything AFTER index i to everything BEFORE index i
				todos = append(todos[:i], todos[i+1:]...)

				c.JSON(200, gin.H{"message": "Todo deleted"})
				store.Save(dbPath, todos) // <--- Use store.Save

				return
			}
		}
		c.JSON(404, gin.H{"message": "Todo not found"})
	})

	r.Run(":8080")
}

func saveTodos() {
	data, _ := json.MarshalIndent(todos, "", "  ") // Pretty print JSON
	os.WriteFile(dbPath, data, 0644)
}
