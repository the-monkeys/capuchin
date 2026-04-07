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

//go:embed migrations/*.sql
var migrations embed.FS

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Config.POSTGRES_HOST,
		config.Config.POSTGRES_USER,
		config.Config.POSTGRES_PASSWORD,
		config.Config.POSTGRES_DB,
		config.Config.POSTGRES_PORT,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("could not connect to database:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
}

// Migrate runs all pending goose migrations embedded in the binary.
// Goose tracks applied versions in the goose_db_version table, making
// repeated calls safe (idempotent).
func Migrate() {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("goose dialect error:", err)
	}
	if err := goose.Up(DB, "migrations"); err != nil {
		log.Fatal("goose migration error:", err)
	}
	log.Println("migrations applied successfully")
}

func CleanupTokens() error {
	_, err := DB.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < NOW()")
	return err
}
