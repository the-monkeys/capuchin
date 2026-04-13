package database

import (
	"capuchin/internal/config"
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

// DB is the shared connection pool. Nil until the background goroutine
// successfully connects for the first time.
var DB *sql.DB

// dbHealthy is 1 when DB is reachable, 0 otherwise.
// Accessed exclusively via sync/atomic.
var dbHealthy int32

const retryInterval = 5 * time.Second

// Connect launches a background goroutine that attempts to open and ping
// Postgres on a fixed interval. It returns immediately without blocking the
// caller — the HTTP server starts before the DB is necessarily ready.
// The backend process never exits due to DB unavailability.
func Connect(cfg config.AppConfig) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	go func() {
		for {
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Printf("database: connection open failed: %v — retry in %s", err, retryInterval)
				atomic.StoreInt32(&dbHealthy, 0)
				time.Sleep(retryInterval)
				continue
			}

			if err := db.Ping(); err != nil {
				log.Printf("database: ping failed: %v — retry in %s", err, retryInterval)
				atomic.StoreInt32(&dbHealthy, 0)
				_ = db.Close()
				time.Sleep(retryInterval)
				continue
			}

			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)

			DB = db
			atomic.StoreInt32(&dbHealthy, 1)
			log.Println("database: connection established")

			// Switch to a periodic health-check ping loop.
			for {
				time.Sleep(retryInterval)
				if err := DB.Ping(); err != nil {
					log.Printf("database: connection lost: %v — reconnecting", err)
					atomic.StoreInt32(&dbHealthy, 0)
					_ = DB.Close()
					DB = nil
					break // fall back to outer reconnect loop
				}
				atomic.StoreInt32(&dbHealthy, 1)
			}
		}
	}()
}

// IsHealthy reports whether the last DB ping succeeded.
func IsHealthy() bool {
	return atomic.LoadInt32(&dbHealthy) == 1
}

// CleanupTokens deletes expired blacklisted tokens.
// Returns an error if the DB is not yet connected.
func CleanupTokens() error {
	if DB == nil {
		return fmt.Errorf("database: not connected")
	}
	_, err := DB.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < NOW()")
	return err
}
