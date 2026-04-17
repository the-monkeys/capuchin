# Implementation Plan: Migration Module Separation

## Overview

Separate all database migration tooling from the `capuchin` backend module into a standalone Go module at `migration/`. The backend is stripped of goose/testcontainers/rapid dependencies and gains a non-blocking `Connect()` with an atomic health flag and a 503 middleware.

## Tasks

- [x] 1. Create the `migration/` Go module scaffold
  - Create `migration/go.mod` with module path `capuchin-migration`, declaring `github.com/pressly/goose/v3`, `github.com/lib/pq`, `github.com/testcontainers/testcontainers-go`, `github.com/testcontainers/testcontainers-go/modules/postgres`, `pgregory.net/rapid`, `github.com/google/uuid`, and `golang.org/x/crypto` as direct dependencies
  - Create `migration/db/embed.go` with `//go:embed migrations/*.sql` pointing at `migration/db/migrations/`
  - Copy `backend/db/migrations/00001_init_schema.sql` → `migration/db/migrations/00001_init_schema.sql`
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3_

- [x] 2. Implement the Migration Runner binary
  - Create `migration/cmd/migrate/main.go` that reads `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT` (default `5432`) directly from `os.Getenv` — no `capuchin` imports
  - Build the DSN string, call `sql.Open`, `db.Ping`, set `goose.SetBaseFS(db.Migrations)`, `goose.SetDialect("postgres")`, `goose.Up`, and exit; log fatal on any error
  - _Requirements: 1.4, 3.1, 3.2, 3.3, 3.4_

- [x] 3. Relocate migration tests to the `migration/` module
  - Create `migration/migrate_test.go` (package `migration_test`) by copying `backend/internal/database/migrate_test.go` verbatim, then update the import of `capuchin/db` → `capuchin-migration/db` and remove the import of `capuchin/internal/database`
  - In `TestP4_AppServerDoesNotMigrate`, replace the `database.DB = db` assignment with a direct `db.Ping()` call (no capuchin import needed — the test only verifies `goose_db_version` does not exist on a fresh DB)
  - Run `go mod tidy` inside `migration/` to populate `go.sum`
  - _Requirements: 5.1, 5.2, 5.3_

  - [x]* 3.1 Verify property tests P1–P6 pass in the migration module
    - **Property 1: Migration file structural invariants** — `TestP1_MigrationFileStructuralInvariants`
    - **Property 2: Migration application round-trip** — `TestP2_MigrationApplicationRoundTrip`
    - **Property 3: Migration idempotency** — `TestP3_MigrationIdempotency`
    - **Property 4: App server does not migrate** — `TestP4_AppServerDoesNotMigrate`
    - **Property 5: Seed runner idempotency** — `TestP5_SeedRunnerIdempotency`
    - **Property 6: Concurrent migration safety** — `TestP6_ConcurrentMigrationSafety`
    - **Validates: Requirements 5.1, 5.2**

- [x] 4. Add `TestInitSQLMatchesMigrationEndState` to the migration module
  - Create `migration/init_sql_test.go` (package `migration_test`) implementing `TestInitSQLMatchesMigrationEndState`
  - Spin up two testcontainers Postgres instances; apply `migration/db/init.sql` to the first via `db.Exec` and `goose.Up` to the second via `migrateDB`; dump both schemas with `pg_dump --schema-only`; fail if the dumps differ
  - _Requirements: 2.2, 2.3_ (init.sql sync check)

  - [x]* 4.1 Verify `TestInitSQLMatchesMigrationEndState` passes
    - **Property 9: init.sql schema matches goose migration end state**
    - **Validates: Requirements 7.1, 7.3**

- [x] 5. Checkpoint — migration module complete
  - Ensure `go build ./...` and `go vet ./...` pass inside `migration/`
  - Ensure all tests pass, ask the user if questions arise.

- [x] 6. Rewrite `backend/internal/database/db.go` — non-blocking Connect with atomic health flag
  - Change signature to `func Connect(cfg config.AppConfig)` (explicit config parameter)
  - Remove the blocking `sql.Open` + `DB.Ping` + `log.Fatal` pattern
  - Add `var dbHealthy int32` (package-level, accessed only via `sync/atomic`)
  - `Connect` launches a background goroutine: loop every 5 seconds — `sql.Open`, `Ping`; on failure log error and set `atomic.StoreInt32(&dbHealthy, 0)`; on success configure pool (`SetMaxOpenConns(25)`, `SetMaxIdleConns(5)`, `SetConnMaxLifetime(5*time.Minute)`), set `atomic.StoreInt32(&dbHealthy, 1)`, log recovery, then switch to a periodic ping loop (same 5-second interval) that sets the flag on each result
  - Add `func IsHealthy() bool { return atomic.LoadInt32(&dbHealthy) == 1 }`
  - Keep `CleanupTokens()` unchanged
  - _Requirements: 6.1, 6.3, 6.4, 6.6, 6.7, 6.8_

  - [x]* 6.1 Write unit test `TestDBHealthFlag_ReflectsConnectionState`
    - Use `database/sql/driver` mock to simulate ping success and failure; assert `IsHealthy()` reflects the last ping result without requiring Docker
    - **Property 7: DB health flag reflects connection state**
    - **Validates: Requirements 6.4, 6.7**

- [x] 7. Create `backend/internal/middleware/db_health.go`
  - Implement `func DBHealthCheck() gin.HandlerFunc` that calls `database.IsHealthy()`; if false, writes `{"error":"database unavailable"}` with status 503 and calls `c.Abort()`; otherwise calls `c.Next()`
  - _Requirements: 6.5_

  - [x]* 7.1 Write unit tests for `DBHealthCheck` middleware
    - `TestDBHealthMiddleware_Returns503WhenUnhealthy`: set `dbHealthy=0`, fire a GET request via `httptest`, assert status 503 and JSON body
    - `TestDBHealthMiddleware_PassesWhenHealthy`: set `dbHealthy=1`, fire a GET request, assert the next handler executes and returns 200
    - **Property 8: 503 returned while unhealthy**
    - **Validates: Requirements 6.5**

- [x] 8. Update `backend/cmd/server/main.go`
  - Change `database.Connect()` → `database.Connect(config.Config)`
  - Register `r.Use(middleware.DBHealthCheck())` immediately after `gin.Default()` and before `routes.SetupRoutes`
  - Add import for `capuchin/internal/middleware`
  - _Requirements: 6.1, 6.3, 6.5_

- [x] 9. Checkpoint — backend changes complete
  - Ensure `go build ./...` and `go vet ./...` pass inside `backend/`
  - Ensure all tests pass, ask the user if questions arise.

- [x] 10. Remove migration artifacts from the backend module
  - Delete `backend/db/embed.go` and `backend/db/migrations/00001_init_schema.sql` (and the `backend/db/` directory if empty)
  - Delete `backend/cmd/migrate/main.go` (and the `backend/cmd/migrate/` directory)
  - Delete `backend/internal/database/migrate_test.go`
  - _Requirements: 2.4, 4.1, 4.2, 4.3, 5.3_

- [x] 11. Clean up `backend/go.mod`
  - Run `go mod tidy` inside `backend/` to remove `github.com/pressly/goose/v3`, `github.com/testcontainers/testcontainers-go`, `github.com/testcontainers/testcontainers-go/modules/postgres`, and `pgregory.net/rapid` from direct and indirect dependencies
  - Verify `github.com/lib/pq` remains as a direct dependency
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [x] 12. Final checkpoint — full separation verified
  - Run `go build ./...` and `go test ./...` inside `backend/` — must pass without Docker or a running Postgres instance
  - Run `go build ./...` and `go vet ./...` inside `migration/` — must compile cleanly
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- No Tasks should be skipped
- Each task references specific requirements for traceability
- The `migration/` module has zero imports from `capuchin` — the two modules are fully independent
- Backend unit tests (tasks 6.1, 7.1) use mock drivers and `httptest` — no Docker required
- Property tests in the migration module (task 3.1) require Docker via testcontainers
