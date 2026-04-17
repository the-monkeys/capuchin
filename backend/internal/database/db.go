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

// dbHealthy is 1 when the last monitor ping succeeded, 0 otherwise.
var dbHealthy atomic.Int32

// degraded is a channel used to signal the monitor to switch to fast polling.
// Buffered so callers never block.
var degraded = make(chan struct{}, 1)

const (
	retryInterval    = 5 * time.Second
	healthyInterval  = 30 * time.Second
	degradedInterval = 5 * time.Second
	pingTimeout      = 2 * time.Second
)

// GetDB returns the active connection pool, or nil if not yet connected.
func GetDB() *sql.DB {
	return db.Load()
}

// IsHealthy reports the cached DB health state set by StartHealthMonitor.
func IsHealthy() bool {
	return dbHealthy.Load() == 1
}

// MarkDegraded signals the monitor to switch to fast polling immediately.
// Safe to call from any goroutine; never blocks.
func MarkDegraded() {
	select {
	case degraded <- struct{}{}:
	default: // already signalled, drop
	}
}

// Connect launches a background goroutine that establishes the DB connection
// and retries on failure. Returns immediately — the server starts without
// waiting for the DB.
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
			dbHealthy.Store(1)
			logger.Info("database", "connection established")
			return
		}
	}()
}

// StartHealthMonitor runs a two-speed ping loop:
//   - healthy:  pings every 30s for observability
//   - degraded: pings every 5s to detect recovery as soon as possible
//
// Switches to degraded mode when a ping fails or MarkDegraded() is called
// (e.g. from a service that hit a query error). Backs off to healthy interval
// once a ping succeeds again.
func StartHealthMonitor() {
	go func() {
		interval := healthyInterval
		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-degraded:
				// A query error was reported — switch to fast polling immediately
				// without waiting for the current timer to fire.
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				interval = degradedInterval
				timer.Reset(interval)

			case <-timer.C:
				ping(interval == degradedInterval)
				if IsHealthy() {
					interval = healthyInterval
				} else {
					interval = degradedInterval
				}
				timer.Reset(interval)
			}
		}
	}()
}

func ping(wasDegraded bool) {
	conn := GetDB()
	if conn == nil {
		dbHealthy.Store(0)
		logger.Warn("database.monitor", "DB not yet connected", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	err := conn.PingContext(ctx)
	cancel()

	if err != nil {
		dbHealthy.Store(0)
		logger.Warn("database.monitor", "DB unreachable", err)
		return
	}

	if wasDegraded {
		logger.Info("database.monitor", "DB recovered")
	}
	dbHealthy.Store(1)
	stats := conn.Stats()
	logger.Info("database.monitor", fmt.Sprintf(
		"healthy — open=%d idle=%d waitCount=%d",
		stats.OpenConnections, stats.Idle, stats.WaitCount,
	))
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
