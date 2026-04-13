package handlers

import (
	"capuchin/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService services.TodoService
}

func NewTodoHandler(svc services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: svc}
}

type addTodoRequest struct {
	Item      string `json:"item" binding:"required"`
	Completed bool   `json:"completed"`
}

type updateTodoRequest struct {
	Item      *string `json:"item"`
	Completed *bool   `json:"completed"`
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	todos, err := h.todoService.GetTodos(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var req addTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.AddTodo(userID, req.Item, req.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add todo"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.todoService.UpdateTodo(userID, id, req.Item, req.Completed)
	if err != nil {
		if err == services.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.todoService.DeleteTodo(userID, id)
	if err != nil {
		if err == services.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
