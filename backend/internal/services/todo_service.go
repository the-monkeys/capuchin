package services

import (
	"capuchin/internal/database"
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
	rows, err := database.DB.Query("SELECT id, item, completed FROM todos WHERE user_id=$1", userID)
	if err != nil {
		return nil, ErrDatabase
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

	if err := rows.Err(); err != nil {
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

	_, err := database.DB.Exec("INSERT INTO todos (id, item, completed, user_id) VALUES ($1, $2, $3, $4)", t.ID, t.Item, t.Completed, t.UserID)
	if err != nil {
		return nil, ErrDatabase
	}
	return t, nil
}

func (s *todoService) UpdateTodo(userID, todoID uuid.UUID, item *string, completed *bool) (*models.Todo, error) {
	if item == nil && completed == nil {
		var t models.Todo
		err := database.DB.QueryRow("SELECT id, item, completed FROM todos WHERE id=$1 AND user_id=$2", todoID, userID).Scan(&t.ID, &t.Item, &t.Completed)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrTodoNotFound
			}
			return nil, ErrDatabase
		}
		return &t, nil
	}

	var t models.Todo
	err := database.DB.QueryRow(`
		UPDATE todos 
		SET item = COALESCE($1, item), 
		    completed = COALESCE($2, completed)
		WHERE id=$3 AND user_id=$4 
		RETURNING id, item, completed`, item, completed, todoID, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, ErrDatabase
	}
	return &t, nil
}

func (s *todoService) DeleteTodo(userID, todoID uuid.UUID) error {
	res, err := database.DB.Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", todoID, userID)
	if err != nil {
		return ErrDatabase
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
