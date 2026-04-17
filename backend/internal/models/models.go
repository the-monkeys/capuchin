package models

type Todo struct {
	ID        int64  `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
	UserID    int64  `json:"-"`
}

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
