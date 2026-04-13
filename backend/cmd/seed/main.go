// seed populates the database with deterministic development data.
// It is idempotent: running it multiple times will not create duplicates.
//
// Usage:
//
//	go run ./cmd/seed
package main

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	id       int64
	email    string
	password string
}

type seedTodo struct {
	id        int64
	userID    int64
	item      string
	completed bool
}

func main() {
	if os.Getenv("APP_ENV") == "production" {
		log.Fatal("seed must not be run in production")
	}

	database.Connect(config.Config)

	// Wait for the background goroutine to establish the DB connection.
	for range 30 {
		if database.IsHealthy() {
			break
		}
		log.Println("seed: database not ready — waiting...")
		time.Sleep(1 * time.Second)
	}
	if !database.IsHealthy() {
		log.Fatal("database not available after 30 seconds")
	}

	users := []seedUser{
		{id: 1, email: "alice@example.com", password: "password123"},
		{id: 2, email: "bob@example.com", password: "password123"},
	}

	todos := []seedTodo{
		{id: 1, userID: 1, item: "Buy groceries", completed: false},
		{id: 2, userID: 1, item: "Read a book", completed: true},
		{id: 3, userID: 2, item: "Go for a run", completed: false},
	}

	seedUsers(users)
	seedTodos(todos)

	log.Println("seed: complete")
}

func seedUsers(users []seedUser) {
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("bcrypt error for %s: %v", u.email, err)
		}

		_, err = database.DB.Exec(`
			INSERT INTO users (id, email, password_hash)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO NOTHING`,
			u.id, u.email, string(hash),
		)
		if err != nil {
			log.Fatalf("failed to seed user %s: %v", u.email, err)
		}
		log.Printf("seed: user inserted: %s", u.email)
	}
}

func seedTodos(todos []seedTodo) {
	for _, t := range todos {
		_, err := database.DB.Exec(`
			INSERT INTO todos (id, item, completed, user_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO NOTHING`,
			t.id, t.item, t.completed, t.userID,
		)
		if err != nil {
			log.Fatalf("failed to seed todo %q: %v", t.item, err)
		}
		log.Printf("seed: todo inserted: %s", t.item)
	}
}
