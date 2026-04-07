// migrate applies pending database migrations and exits.
// Run this as a one-off job (CI/CD step or init container) before deploying
// app server instances. The app server has no migration logic.
//
// Usage:
//
//	go run ./cmd/migrate
package main

import (
	capuchindb "capuchin/db"
	"capuchin/internal/config"
	"capuchin/internal/database"
	"log"

	"github.com/pressly/goose/v3"
)

func main() {
	_ = config.Config
	database.Connect()

	goose.SetBaseFS(capuchindb.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("goose dialect error:", err)
	}
	if err := goose.Up(database.DB, "migrations"); err != nil {
		log.Fatal("goose migration error:", err)
	}
	log.Println("migrations applied successfully")
}
