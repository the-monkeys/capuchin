package models

import "github.com/google/uuid"

type Todo struct {
	ID        uuid.UUID `json:"id"`
	Item      string    `json:"item"`
	Completed bool      `json:"completed"`
	UserID    uuid.UUID `json:"-"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
}
