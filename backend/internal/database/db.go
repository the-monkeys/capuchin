package database

import (
	"capuchin/internal/config"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var migrations embed.FS

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Config.POSTGRES_HOST,
		config.Config.POSTGRES_USER,
		config.Config.POSTGRES_PASSWORD,
		config.Config.POSTGRES_DB,
		config.Config.POSTGRES_PORT,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// Conservative pool settings avoid exhausting DB connections in small deployments.
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	// Recycling connections helps recover from stale network state over long uptimes.
	DB.SetConnMaxLifetime(5 * time.Minute)
}

func Migrate() {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("goose dialect error:", err)
	}
	if err := goose.Up(DB, "migrations"); err != nil {
		log.Fatal("goose migration error:", err)
	}
	log.Println("Database migrations applied successfully.")
}

func CleanupTokens() error {
	// Expired tokens can be dropped because JWT expiration already invalidates them.
	_, err := DB.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < $1", time.Now())
	return err
}
