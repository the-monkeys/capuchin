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

// db is the shared connection pool. Access via GetDB().
// Uses atomic.Pointer to avoid data races on connect/reconnect.
var db atomic.Pointer[sql.DB]

// dbHealthy is 1 when the DB connection is established, 0 otherwise.
var dbHealthy int32

const retryInterval = 5 * time.Second

// GetDB returns the active connection pool, or nil if not yet connected.
func GetDB() *sql.DB {
	return db.Load()
}

// IsHealthy reports whether the DB is currently reachable.
func IsHealthy() bool {
	return atomic.LoadInt32(&dbHealthy) == 1
}

// Connect launches a background goroutine that establishes the DB connection
// and retries on failure. Returns immediately — the server starts without
// waiting for the DB. Once connected, sql.DB manages the pool internally;
// no polling loop is needed.
func Connect(cfg config.AppConfig) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_DB,
		cfg.POSTGRES_PORT,
	)

	go func() {
		for {
			conn, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Printf("database: open failed: %v — retry in %s", err, retryInterval)
				time.Sleep(retryInterval)
				continue
			}

			if err := conn.Ping(); err != nil {
				log.Printf("database: ping failed: %v — retry in %s", err, retryInterval)
				_ = conn.Close()
				atomic.StoreInt32(&dbHealthy, 0)
				time.Sleep(retryInterval)
				continue
			}

			conn.SetMaxOpenConns(25)
			conn.SetMaxIdleConns(5)
			conn.SetConnMaxLifetime(5 * time.Minute)

			db.Store(conn)
			atomic.StoreInt32(&dbHealthy, 1)
			log.Println("database: connection established")
			return
		}
	}()
}

// CleanupTokens deletes expired blacklisted tokens.
func CleanupTokens() error {
	conn := GetDB()
	if conn == nil {
		return fmt.Errorf("database: not connected")
	}
	_, err := conn.Exec("DELETE FROM blacklisted_tokens WHERE expired_at < NOW()")
	return err
}
