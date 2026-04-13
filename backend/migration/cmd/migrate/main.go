// migrate applies pending database migrations and exits.
// Run this as a one-off job (CI/CD step or init container) before deploying
// app server instances.
//
// Usage:
//
//	go run ./cmd/migrate
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	migrationdb "capuchin-migration/db"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		log.Fatal("POSTGRES_HOST environment variable is required")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		log.Fatal("POSTGRES_USER environment variable is required")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		log.Fatal("POSTGRES_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		log.Fatal("POSTGRES_DB environment variable is required")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database connection:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	goose.SetBaseFS(migrationdb.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("goose dialect error:", err)
	}

	if err := goose.Up(db, "versions"); err != nil {
		log.Fatal("goose migration error:", err)
	}

	log.Println("migrate: migrations applied")
}
