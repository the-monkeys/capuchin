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
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Fixed UUIDs keep seed data stable across runs so foreign keys stay consistent.
var (
	user1ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

type seedUser struct {
	id       uuid.UUID
	email    string
	password string
}

type seedTodo struct {
	id        uuid.UUID
	userID    uuid.UUID
	item      string
	completed bool
}

func main() {
	// config.init() runs automatically on import; Connect() needs explicit call.
	database.Connect(config.Config)

	// Wait for the background goroutine to establish the DB connection.
	for i := 0; i < 30; i++ {
		if database.IsHealthy() {
			break
		}
		log.Println("waiting for database connection...")
		time.Sleep(1 * time.Second)
	}
	if !database.IsHealthy() {
		log.Fatal("database not available after 30 seconds")
	}

	users := []seedUser{
		{id: user1ID, email: "alice@example.com", password: "password123"},
		{id: user2ID, email: "bob@example.com", password: "password123"},
	}

	todos := []seedTodo{
		{id: uuid.MustParse("00000000-0000-0000-0001-000000000001"), userID: user1ID, item: "Buy groceries", completed: false},
		{id: uuid.MustParse("00000000-0000-0000-0001-000000000002"), userID: user1ID, item: "Read a book", completed: true},
		{id: uuid.MustParse("00000000-0000-0000-0001-000000000003"), userID: user2ID, item: "Go for a run", completed: false},
	}

	seedUsers(users)
	seedTodos(todos)

	log.Println("seed complete")
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
		log.Printf("seeded user: %s", u.email)
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
		log.Printf("seeded todo: %s", t.item)
	}
}
