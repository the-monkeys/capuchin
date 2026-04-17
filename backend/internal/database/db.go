package database

import (
	"capuchin/internal/config"
	"capuchin/internal/logger"
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

// db is the shared connection pool. Access via GetDB().
// Uses atomic.Pointer to avoid data races on connect/reconnect.
var db atomic.Pointer[sql.DB]

const (
	retryInterval   = 5 * time.Second
	monitorInterval = 30 * time.Second
	pingTimeout     = 2 * time.Second
)

// GetDB returns the active connection pool, or nil if not yet connected.
func GetDB() *sql.DB {
	return db.Load()
}

// Connect launches a background goroutine that establishes the DB connection
// and retries on failure. Returns immediately — the server starts without
// waiting for the DB. Once connected, sql.DB manages the pool internally.
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
				logger.Error("database", "open failed", err)
				time.Sleep(retryInterval)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err = conn.PingContext(ctx)
			cancel()

			if err != nil {
				logger.Error("database", "ping failed during connect", err)
				_ = conn.Close()
				time.Sleep(retryInterval)
				continue
			}

			conn.SetMaxOpenConns(25)
			conn.SetMaxIdleConns(5)
			conn.SetConnMaxLifetime(5 * time.Minute)

			db.Store(conn)
			logger.Info("database", "connection established")
			return
		}
	}()
}

// StartHealthMonitor pings the DB every 30s and logs a structured warning
// when unreachable. Intended for observability only — does not gate requests.
// Call after Connect().
func StartHealthMonitor() {
	go func() {
		ticker := time.NewTicker(monitorInterval)
		defer ticker.Stop()
		for range ticker.C {
			conn := GetDB()
			if conn == nil {
				logger.Warn("database.monitor", "DB not yet connected", nil)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err := conn.PingContext(ctx)
			cancel()

			if err != nil {
				logger.Warn("database.monitor", "DB unreachable", err)
			} else {
				stats := conn.Stats()
				logger.Info("database.monitor", fmt.Sprintf(
					"healthy — open=%d idle=%d waitCount=%d",
					stats.OpenConnections, stats.Idle, stats.WaitCount,
				))
			}
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
