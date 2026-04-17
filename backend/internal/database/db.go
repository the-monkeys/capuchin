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
	connectRetryInterval = 5 * time.Second
	healthyPingInterval  = 30 * time.Second
	degradedPingInterval = 5 * time.Second
	pingTimeout          = 2 * time.Second
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

// Connect opens the connection pool once and launches a background goroutine
// that pings until ready, retrying on failure. Returns immediately.
func Connect(cfg config.AppConfig) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_DB,
		cfg.POSTGRES_PORT,
	)

	// sql.Open only validates the DSN — allocate the pool once outside the retry loop.
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("database", "failed to open connection pool", err)
		return
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err := conn.PingContext(ctx)
			cancel()

			if err != nil {
				logger.Error("database", "ping failed, retrying", err)
				time.Sleep(connectRetryInterval)
				continue
			}

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
// Switches to degraded mode when a ping fails or MarkDegraded() is called.
// Backs off to healthy interval once a ping succeeds again.
func StartHealthMonitor() {
	go func() {
		isDegraded := false
		timer := time.NewTimer(healthyPingInterval)
		defer timer.Stop()

		for {
			select {
			case <-degraded:
				// A query error was reported — switch to fast polling immediately
				// without waiting for the current timer to fire.
				if !isDegraded {
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					isDegraded = true
					timer.Reset(degradedPingInterval)
				}

			case <-timer.C:
				pingMonitor(isDegraded)
				if IsHealthy() {
					isDegraded = false
					timer.Reset(healthyPingInterval)
				} else {
					isDegraded = true
					timer.Reset(degradedPingInterval)
				}
			}
		}
	}()
}

func pingMonitor(isDegraded bool) {
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

	if isDegraded {
		// Only log recovery and stats when coming back from a degraded state.
		stats := conn.Stats()
		logger.Info("database.monitor", fmt.Sprintf(
			"DB recovered — open=%d idle=%d waitCount=%d",
			stats.OpenConnections, stats.Idle, stats.WaitCount,
		))
	}
	dbHealthy.Store(1)
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
