package todo

import (
	"capuchin/internal/database"
	"capuchin/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTodos(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	rows, err := database.DB.Query("SELECT id, item, completed FROM todos WHERE user_id=$1", userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
			continue
		}
		todos = append(todos, t)
	}
	c.JSON(200, todos)
}

func AddTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var t models.Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	t.ID = uuid.New()
	t.UserID = userID

	_, err := database.DB.Exec("INSERT INTO todos (id, item, completed, user_id) VALUES ($1, $2, $3, $4)", t.ID, t.Item, t.Completed, t.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, t)
}

func ToggleTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var t models.Todo
	// Toggle and return new state
	err = database.DB.QueryRow(`
		UPDATE todos SET completed = NOT completed 
		WHERE id=$1 AND user_id=$2 
		RETURNING id, item, completed`, id, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(200, t)
}

func EditTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var t models.Todo
	err = database.DB.QueryRow(`
		UPDATE todos SET item=$1 
		WHERE id=$2 AND user_id=$3 
		RETURNING id, item, completed`, req.Item, id, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(200, t)
}

func DeleteTodo(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := database.DB.Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete todo"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// This means the todo didn't exist or didn't belong to the user
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Todo deleted successfully"})
}
