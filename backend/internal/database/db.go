package database

import (
	"capuchin/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	dbMaxStartupAttempts = 5
	dbStartupBaseDelay   = 2 * time.Second
	dbStartupMaxDelay    = 30 * time.Second
)

func Connect() error {
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
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Conservative pool settings avoid exhausting DB connections in small deployments.
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	// Recycling connections helps recover from stale network state over long uptimes.
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = pingWithRetry(); err != nil {
		DB.Close()
		return fmt.Errorf("database unavailable after %d attempts: %w", dbMaxStartupAttempts, err)
	}

	log.Println("Database connection established")
	return nil
}

func pingWithRetry() error {
	delay := dbStartupBaseDelay
	for attempt := 1; attempt <= dbMaxStartupAttempts; attempt++ {
		if err := DB.Ping(); err == nil {
			return nil
		} else {
			log.Printf("Database ping failed (attempt %d/%d): %v", attempt, dbMaxStartupAttempts, err)
		}
		if attempt < dbMaxStartupAttempts {
			log.Printf("Retrying in %s...", delay)
			time.Sleep(delay)
			delay *= 2
			if delay > dbStartupMaxDelay {
				delay = dbStartupMaxDelay
			}
		}
	}
	return fmt.Errorf("database unreachable after %d attempts", dbMaxStartupAttempts)
}

func InitSchema() {
	log.Println("Database connection initialized. Assuming schema is already present.")
}

func CleanupTokens() error {
	// Expired tokens can be dropped because JWT expiration already invalidates them.
	_, err := DB.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < $1", time.Now())
	return err
}
