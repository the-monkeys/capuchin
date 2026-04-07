// migrate applies pending database migrations and exits.
// Run this as a one-off job before deploying new app instances.
//
// Usage:
//
//	go run ./cmd/migrate
package main

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"log"
)

func main() {
	_ = config.Config
	database.Connect()
	database.Migrate()
	log.Println("done")
}
