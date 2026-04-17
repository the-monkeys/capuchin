package services

import (
	"capuchin/internal/database"
	"capuchin/internal/logger"
	"capuchin/internal/models"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

type TodoService interface {
	GetTodos(userID uuid.UUID) ([]models.Todo, error)
	AddTodo(userID uuid.UUID, item string, completed bool) (*models.Todo, error)
	UpdateTodo(userID, todoID uuid.UUID, item *string, completed *bool) (*models.Todo, error)
	DeleteTodo(userID, todoID uuid.UUID) error
}

type todoService struct{}

func NewTodoService() TodoService {
	return &todoService{}
}

func (s *todoService) GetTodos(userID uuid.UUID) ([]models.Todo, error) {
	// Scope every read by user_id so one user can never read another user's todos.
	rows, err := database.GetDB().Query("SELECT id, item, completed FROM todos WHERE user_id=$1", userID)
	if err != nil {
		logger.Error("todo.get", "query failed", err)
		return nil, ErrDatabase
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
			// Skip malformed rows instead of failing the whole response for a single bad record.
			continue
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		logger.Error("todo.get", "rows iteration failed", err)
		return nil, ErrDatabase
	}

	return todos, nil
}

func (s *todoService) AddTodo(userID uuid.UUID, item string, completed bool) (*models.Todo, error) {
	t := &models.Todo{
		ID:        uuid.New(),
		UserID:    userID,
		Item:      item,
		Completed: completed,
	}

	_, err := database.GetDB().Exec("INSERT INTO todos (id, item, completed, user_id) VALUES ($1, $2, $3, $4)", t.ID, t.Item, t.Completed, t.UserID)
	if err != nil {
		logger.Error("todo.add", "insert failed", err)
		return nil, ErrDatabase
	}
	return t, nil
}

func (s *todoService) UpdateTodo(userID, todoID uuid.UUID, item *string, completed *bool) (*models.Todo, error) {
	if item == nil && completed == nil {
		// Empty PATCH requests are treated as a read to keep the endpoint idempotent.
		var t models.Todo
		err := database.GetDB().QueryRow("SELECT id, item, completed FROM todos WHERE id=$1 AND user_id=$2", todoID, userID).Scan(&t.ID, &t.Item, &t.Completed)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrTodoNotFound
			}
			logger.Error("todo.update", "read-only fetch failed", err)
			return nil, ErrDatabase
		}
		return &t, nil
	}

	var t models.Todo
	err := database.GetDB().QueryRow(`
		UPDATE todos 
		-- COALESCE preserves existing values when fields are omitted from PATCH payloads.
		SET item = COALESCE($1, item), 
		    completed = COALESCE($2, completed)
		WHERE id=$3 AND user_id=$4 
		RETURNING id, item, completed`, item, completed, todoID, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		logger.Error("todo.update", "update query failed", err)
		return nil, ErrDatabase
	}
	return &t, nil
}

func (s *todoService) DeleteTodo(userID, todoID uuid.UUID) error {
	res, err := database.GetDB().Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", todoID, userID)
	if err != nil {
		logger.Error("todo.delete", "delete query failed", err)
		return ErrDatabase
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// Distinguish "not found" from successful deletion for better API semantics.
		return ErrTodoNotFound
	}

	return nil
}
