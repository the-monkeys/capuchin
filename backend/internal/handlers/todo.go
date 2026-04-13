package handlers

import (
	"capuchin/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService services.TodoService
}

func NewTodoHandler(svc services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: svc}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	todos, err := h.todoService.GetTodos(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get todos"})
		return
	}
	c.JSON(200, todos)
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req struct {
		Item      string `json:"item" binding:"required"`
		Completed bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.AddTodo(userID, req.Item, req.Completed)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add todo"})
		return
	}
	c.JSON(200, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Item      *string `json:"item"`
		Completed *bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.UpdateTodo(userID, id, req.Item, req.Completed)
	if err != nil {
		if err == services.ErrTodoNotFound {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(200, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	err = h.todoService.DeleteTodo(userID, id)
	if err != nil {
		if err == services.ErrTodoNotFound {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(200, gin.H{"message": "Todo deleted successfully"})
}
