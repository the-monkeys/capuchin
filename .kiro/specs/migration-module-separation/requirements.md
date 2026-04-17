# Requirements Document

## Introduction

This feature separates all database migration tooling from the main backend application (`capuchin` module) into a standalone Go module living at `migration/` in the repository root. The migration module owns the goose runner, embedded SQL files, migration binary, and migration tests. The backend module is stripped of all migration-related dependencies (goose, testcontainers, rapid) and is entirely unaware of migration state. The backend process lifetime is fully independent of Postgres: it starts immediately, attempts to connect in the background, and continues running regardless of whether Postgres is reachable — at startup or at any point during runtime. Migrations are release-gated: the migration pipeline runs only during a release, and a failed migration blocks deployment of the new backend version.

## Glossary

- **Migration_Module**: The standalone Go module rooted at `migration/` with its own `go.mod`, responsible for all schema migration concerns.
- **Backend_Module**: The existing Go module rooted at `backend/` (`module capuchin`), responsible only for serving the application.
- **Migration_Runner**: The binary entrypoint (`migration/cmd/migrate/main.go`) that applies pending goose migrations and exits.
- **Migration_FS**: The embedded `embed.FS` in the Migration_Module that holds all `*.sql` goose migration files.
- **Goose**: The third-party migration library (`github.com/pressly/goose/v3`) used by the Migration_Module.
- **Testcontainers**: The third-party library (`github.com/testcontainers/testcontainers-go`) used by migration tests to spin up ephemeral Postgres instances.
- **Release_Pipeline**: The CI/CD workflow that runs the Migration_Runner before deploying a new backend version.
- **Backend_Server**: The binary entrypoint (`backend/cmd/server/main.go`) that serves the application; it calls only `database.Connect()`.

---

## Requirements

### Requirement 1: Migration Module Structure

**User Story:** As a platform engineer, I want the migration tooling to live in its own Go module, so that the backend build is not coupled to migration dependencies.

#### Acceptance Criteria

1. THE Migration_Module SHALL reside at `migration/` in the repository root with its own `go.mod` declaring a module path distinct from `capuchin`.
2. THE Migration_Module SHALL declare `github.com/pressly/goose/v3`, `github.com/lib/pq`, and `github.com/testcontainers/testcontainers-go` as direct dependencies in its `go.mod`.
3. THE Migration_Module SHALL embed all `*.sql` goose migration files via a Go `embed.FS` declared in a dedicated package (e.g., `migration/db/embed.go`).
4. THE Migration_Module SHALL contain the Migration_Runner binary at `migration/cmd/migrate/main.go`.

---

### Requirement 2: Migration SQL and Embedded FS Relocation

**User Story:** As a platform engineer, I want all SQL migration files and the embed directive to live exclusively in the Migration_Module, so that the backend has no compile-time dependency on migration assets.

#### Acceptance Criteria

1. THE Migration_Module SHALL contain all goose SQL migration files previously located at `backend/db/migrations/`.
2. THE Migration_FS SHALL embed the SQL files using `//go:embed migrations/*.sql` within the Migration_Module.
3. WHEN the Migration_Runner is compiled, THE Migration_Module SHALL resolve all embedded SQL files without referencing any path inside `backend/`.
4. THE Backend_Module SHALL contain no `embed.FS` referencing migration SQL files after the relocation.

---

### Requirement 3: Migration Runner Relocation

**User Story:** As a platform engineer, I want the migration binary to live in the Migration_Module, so that it can be built and run independently of the backend.

#### Acceptance Criteria

1. THE Migration_Runner SHALL read database connection parameters from environment variables (host, user, password, dbname, port).
2. WHEN all required environment variables are present and Postgres is reachable, THE Migration_Runner SHALL apply all pending goose migrations and exit with code 0.
3. IF a migration step fails, THEN THE Migration_Runner SHALL log the error and exit with a non-zero exit code.
4. THE Migration_Runner SHALL NOT import any package from the `capuchin` module.

---

### Requirement 4: Backend Module Dependency Cleanup

**User Story:** As a backend developer, I want the backend `go.mod` to be free of migration-only dependencies, so that backend builds are faster and the dependency surface is smaller.

#### Acceptance Criteria

1. THE Backend_Module's `go.mod` SHALL NOT list `github.com/pressly/goose/v3` as a direct or indirect dependency after the separation.
2. THE Backend_Module's `go.mod` SHALL NOT list `github.com/testcontainers/testcontainers-go` or `github.com/testcontainers/testcontainers-go/modules/postgres` as direct or indirect dependencies after the separation.
3. THE Backend_Module's `go.mod` SHALL NOT list `pgregory.net/rapid` as a direct or indirect dependency after the separation, provided no remaining backend test uses it.
4. THE Backend_Module SHALL continue to declare `github.com/lib/pq` as a direct dependency because `database.Connect()` requires the Postgres driver.

---

### Requirement 5: Migration Tests Relocation

**User Story:** As a platform engineer, I want all migration-specific tests to live in the Migration_Module, so that the backend test suite has no dependency on Testcontainers or Goose.

#### Acceptance Criteria

1. THE Migration_Module SHALL contain all property-based and integration tests previously located at `backend/internal/database/migrate_test.go`.
2. WHEN the migration test suite is executed, THE Migration_Module SHALL run all six property tests (structural invariants, round-trip, idempotency, app-server-does-not-migrate, seed idempotency, concurrent safety) without importing any package from `capuchin`.
3. THE Backend_Module's test suite SHALL NOT import `github.com/pressly/goose/v3`, `github.com/testcontainers/testcontainers-go`, or `pgregory.net/rapid` after the relocation.
4. WHEN `go test ./...` is run inside `backend/`, THE Backend_Module SHALL compile and pass without requiring Docker or a running Postgres instance.

---

### Requirement 6: Backend Server Startup and DB Health Independence

**User Story:** As a backend developer, I want the backend server to run independently of Postgres availability at all times, so that neither a startup-time nor a runtime database outage ever terminates the server process, and the system recovers automatically when Postgres becomes reachable again.

#### Acceptance Criteria

1. WHEN the Backend_Server starts, THE Backend_Module SHALL call only `database.Connect()` and SHALL NOT invoke any goose or migration function.
2. THE Backend_Module SHALL NOT perform any schema creation or alteration at runtime.
3. WHEN the Backend_Server starts, THE Backend_Module SHALL immediately begin attempting to connect to Postgres in the background without blocking the HTTP server from starting.
4. IF Postgres is unreachable at any point — whether at startup or during runtime — THEN THE Backend_Module SHALL log the connection error, mark the database as unhealthy, and continue running without exiting.
5. WHILE the database is marked unhealthy, THE Backend_Module SHALL return an HTTP 503 Service Unavailable response to any request that requires database access.
6. WHILE the database is marked unhealthy, THE Backend_Module SHALL continuously retry the Postgres connection at a fixed interval.
7. WHEN a retry attempt succeeds, THE Backend_Module SHALL mark the database as healthy and resume normal request handling without requiring a process restart.
8. THE Backend_Server process SHALL never exit due to Postgres being unreachable — the Backend_Server process lifetime is fully independent of Postgres availability.

---

### Requirement 7: Release Pipeline Gating

**User Story:** As a release engineer, I want the migration step to gate the backend deployment, so that a failed migration prevents a broken schema from reaching production.

#### Acceptance Criteria

1. WHEN a new backend version is released, THE Release_Pipeline SHALL execute the Migration_Runner before deploying the Backend_Server.
2. IF the Migration_Runner exits with a non-zero exit code, THEN THE Release_Pipeline SHALL halt and SHALL NOT deploy the new Backend_Server version.
3. WHEN the Migration_Runner exits with code 0, THE Release_Pipeline SHALL proceed to deploy the Backend_Server.
4. THE Release_Pipeline SHALL NOT execute the Migration_Runner on every backend restart — only on release events.
5. THE Backend_Server SHALL NOT require migrations to have been applied as a precondition for startup — migration gating is a release concern only.
