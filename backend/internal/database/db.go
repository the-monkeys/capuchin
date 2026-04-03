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
	maxRetries    = 5
	retryInterval = 5 * time.Second
)

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

	// Conservative pool settings avoid exhausting DB connections in small deployments.
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	// Recycling connections helps recover from stale network state over long uptimes.
	DB.SetConnMaxLifetime(5 * time.Minute)

	for i := 1; i <= maxRetries; i++ {
		if err = DB.Ping(); err == nil {
			if i > 1 {
				log.Printf("Database connected successfully on attempt %d.", i)
			} else {
				log.Println("Database connected successfully.")
			}
			return
		}

		if i < maxRetries {
			log.Printf("DB connection attempt %d/%d failed: %v. Retrying in %v...", i, maxRetries, err, retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Fatalf("could not connect to database after %d attempts: %v — shutting down", maxRetries, err)
}

func InitSchema() {
	// Schema is assumed to be pre-initialized (e.g., via CI/CD pipelines).
	log.Println("Database connection initialized. Assuming schema is already present.")
}

func CleanupTokens() error {
	// Expired tokens can be dropped because JWT expiration already invalidates them.
	_, err := DB.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < $1", time.Now())
	return err
}
