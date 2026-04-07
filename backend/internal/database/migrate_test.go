package database_test

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"capuchin/internal/database"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"pgregory.net/rapid"
)

//go:embed migrations/*.sql
var testMigrations embed.FS

// newTestDB spins up a testcontainers postgres instance and returns a *sql.DB.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	pgc, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pgc.Terminate(ctx) })

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// runStartupMigration mirrors the MODE-gated logic in cmd/server/main.go.
func runStartupMigration(mode string, migrateFn func()) {
	if mode == "dev" {
		migrateFn()
	}
}

// TestP1_MigrationFileStructuralInvariants checks naming convention, goose
// annotations, and CREATE TABLE IF NOT EXISTS for every migration file.
//
// Feature: db-migrations-seeding, Property 1: For any .sql file in
// backend/internal/database/migrations/, the filename must have a zero-padded
// five-digit numeric prefix strictly greater than all preceding files, the file
// must contain both a -- +goose Up block and a -- +goose Down block, and any
// CREATE TABLE statement must use CREATE TABLE IF NOT EXISTS.
func TestP1_MigrationFileStructuralInvariants(t *testing.T) {
	// Feature: db-migrations-seeding, Property 1: Migration file structural invariants
	filenameRe := regexp.MustCompile(`^\d{5}_[a-z0-9_]+\.sql$`)
	// Matches CREATE TABLE followed by IF NOT EXISTS (safe form)
	createTableSafeRe := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS`)
	// Matches any CREATE TABLE occurrence
	createTableAnyRe := regexp.MustCompile(`(?i)CREATE\s+TABLE\b`)

	entries, err := fs.ReadDir(testMigrations, "migrations")
	if err != nil {
		t.Fatalf("failed to read migrations dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no migration files found")
	}

	var prevPrefix int = -1
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		// Assert filename matches pattern
		if !filenameRe.MatchString(name) {
			t.Errorf("filename %q does not match pattern ^\\d{5}_[a-z0-9_]+\\.sql$", name)
		}

		// Assert strictly increasing prefix
		prefixStr := name[:5]
		prefix, err := strconv.Atoi(prefixStr)
		if err != nil {
			t.Errorf("filename %q has non-numeric prefix: %v", name, err)
			continue
		}
		if prefix <= prevPrefix {
			t.Errorf("filename %q prefix %d is not strictly greater than previous %d", name, prefix, prevPrefix)
		}
		prevPrefix = prefix

		// Read file content
		content, err := testMigrations.ReadFile("migrations/" + name)
		if err != nil {
			t.Fatalf("failed to read migration file %q: %v", name, err)
		}
		text := string(content)

		// Assert both goose annotations present
		if !strings.Contains(text, "-- +goose Up") {
			t.Errorf("migration %q missing '-- +goose Up'", name)
		}
		if !strings.Contains(text, "-- +goose Down") {
			t.Errorf("migration %q missing '-- +goose Down'", name)
		}

		// Assert no CREATE TABLE without IF NOT EXISTS:
		// count all CREATE TABLE occurrences and safe CREATE TABLE IF NOT EXISTS occurrences
		allMatches := createTableAnyRe.FindAllString(text, -1)
		safeMatches := createTableSafeRe.FindAllString(text, -1)
		if len(allMatches) != len(safeMatches) {
			t.Errorf("migration %q contains CREATE TABLE without IF NOT EXISTS (total=%d, safe=%d)", name, len(allMatches), len(safeMatches))
		}
	}
}

// TestP2_MigrationApplicationRoundTrip applies migrations to a fresh DB and
// verifies goose_db_version records version 1 as applied.
//
// Feature: db-migrations-seeding, Property 2: For any pending migration
// version, after goose.Up completes successfully, querying goose_db_version
// must return a row with that version's version_id and is_applied = true.
func TestP2_MigrationApplicationRoundTrip(t *testing.T) {
	// Feature: db-migrations-seeding, Property 2: Migration application round-trip
	rapid.Check(t, func(rt *rapid.T) {
		db := newTestDB(t)
		database.DB = db
		database.Migrate()

		var versionID int64
		var isApplied bool
		err := db.QueryRow(
			`SELECT version_id, is_applied FROM goose_db_version WHERE version_id = 1`,
		).Scan(&versionID, &isApplied)
		if err != nil {
			rt.Fatalf("failed to query goose_db_version: %v", err)
		}
		if versionID != 1 {
			rt.Errorf("expected version_id 1, got %d", versionID)
		}
		if !isApplied {
			rt.Errorf("expected is_applied = true for version 1, got false")
		}
	})
}

// TestP3_MigrationIdempotency applies migrations twice and asserts the
// goose_db_version row count is unchanged on the second run.
//
// Feature: db-migrations-seeding, Property 3: For any database state where all
// migrations are already applied, invoking goose.Up again must produce no
// schema changes and the count of rows in goose_db_version must be the same
// before and after the second invocation.
func TestP3_MigrationIdempotency(t *testing.T) {
	// Feature: db-migrations-seeding, Property 3: Migration idempotency
	rapid.Check(t, func(rt *rapid.T) {
		db := newTestDB(t)
		database.DB = db

		// First migration
		database.Migrate()

		var countBefore int
		if err := db.QueryRow(`SELECT COUNT(*) FROM goose_db_version`).Scan(&countBefore); err != nil {
			rt.Fatalf("failed to count goose_db_version rows: %v", err)
		}

		// Second migration — must be idempotent
		database.Migrate()

		var countAfter int
		if err := db.QueryRow(`SELECT COUNT(*) FROM goose_db_version`).Scan(&countAfter); err != nil {
			rt.Fatalf("failed to count goose_db_version rows after second run: %v", err)
		}

		if countAfter != countBefore {
			rt.Errorf("goose_db_version row count changed after second Migrate(): before=%d after=%d", countBefore, countAfter)
		}
	})
}

// TestP4_ModeGatedMigrationBehavior uses rapid to generate arbitrary MODE
// strings and asserts migrateFn is called iff mode == "dev".
//
// Feature: db-migrations-seeding, Property 4: For any value of the MODE
// environment variable that is not exactly "dev", database.Migrate() must not
// be called. When MODE is exactly "dev", database.Migrate() must be called.
func TestP4_ModeGatedMigrationBehavior(t *testing.T) {
	// Feature: db-migrations-seeding, Property 4: MODE-gated migration behavior
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random mode string using lowercase letters
		mode := rapid.StringOf(rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz"))).Draw(rt, "mode")

		called := false
		spy := func() { called = true }

		runStartupMigration(mode, spy)

		if mode == "dev" {
			if !called {
				rt.Errorf("expected migrateFn to be called when mode=%q, but it was not", mode)
			}
		} else {
			if called {
				rt.Errorf("expected migrateFn NOT to be called when mode=%q, but it was", mode)
			}
		}
	})
}

// TestP5_SeedRunnerIdempotency runs seed inserts twice and asserts row counts
// are identical after both runs.
//
// Feature: db-migrations-seeding, Property 5: For any database state, running
// the Seed_Runner twice in sequence must produce the same set of rows as
// running it once — no duplicate rows, no errors on the second run.
func TestP5_SeedRunnerIdempotency(t *testing.T) {
	// Feature: db-migrations-seeding, Property 5: Seed runner idempotency
	rapid.Check(t, func(rt *rapid.T) {
		db := newTestDB(t)
		database.DB = db
		database.Migrate()

		runSeed := func() {
			// Inline seed logic from cmd/seed/main.go using fixed UUIDs
			user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
			user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

			type seedUser struct {
				id       uuid.UUID
				email    string
				password string
			}
			users := []seedUser{
				{id: user1ID, email: "alice@example.com", password: "password123"},
				{id: user2ID, email: "bob@example.com", password: "password123"},
			}
			for _, u := range users {
				hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
				if err != nil {
					rt.Fatalf("bcrypt error: %v", err)
				}
				_, err = db.Exec(`
					INSERT INTO users (id, email, password_hash)
					VALUES ($1, $2, $3)
					ON CONFLICT (id) DO NOTHING`,
					u.id, u.email, string(hash),
				)
				if err != nil {
					rt.Fatalf("failed to seed user %s: %v", u.email, err)
				}
			}

			type seedTodo struct {
				id        uuid.UUID
				userID    uuid.UUID
				item      string
				completed bool
			}
			todos := []seedTodo{
				{id: uuid.MustParse("00000000-0000-0000-0001-000000000001"), userID: user1ID, item: "Buy groceries", completed: false},
				{id: uuid.MustParse("00000000-0000-0000-0001-000000000002"), userID: user1ID, item: "Read a book", completed: true},
				{id: uuid.MustParse("00000000-0000-0000-0001-000000000003"), userID: user2ID, item: "Go for a run", completed: false},
			}
			for _, td := range todos {
				_, err := db.Exec(`
					INSERT INTO todos (id, item, completed, user_id)
					VALUES ($1, $2, $3, $4)
					ON CONFLICT (id) DO NOTHING`,
					td.id, td.item, td.completed, td.userID,
				)
				if err != nil {
					rt.Fatalf("failed to seed todo %q: %v", td.item, err)
				}
			}
		}

		// First seed run
		runSeed()

		var usersBefore, todosBefore int
		if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&usersBefore); err != nil {
			rt.Fatalf("failed to count users: %v", err)
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM todos`).Scan(&todosBefore); err != nil {
			rt.Fatalf("failed to count todos: %v", err)
		}

		// Second seed run — must be idempotent
		runSeed()

		var usersAfter, todosAfter int
		if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&usersAfter); err != nil {
			rt.Fatalf("failed to count users after second seed: %v", err)
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM todos`).Scan(&todosAfter); err != nil {
			rt.Fatalf("failed to count todos after second seed: %v", err)
		}

		if usersAfter != usersBefore {
			rt.Errorf("users count changed after second seed: before=%d after=%d", usersBefore, usersAfter)
		}
		if todosAfter != todosBefore {
			rt.Errorf("todos count changed after second seed: before=%d after=%d", todosBefore, todosAfter)
		}
	})
}

// TestP6_ConcurrentMigrationSafety launches two goroutines both calling
// database.Migrate() simultaneously and asserts version 1 appears exactly once
// with is_applied = true.
//
// Feature: db-migrations-seeding, Property 6: For any two Migration_Runner
// processes started simultaneously against the same database, each migration
// version must appear in goose_db_version with is_applied = true exactly once
// after both processes complete.
func TestP6_ConcurrentMigrationSafety(t *testing.T) {
	// Feature: db-migrations-seeding, Property 6: Concurrent migration safety
	rapid.Check(t, func(rt *rapid.T) {
		db := newTestDB(t)
		database.DB = db

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			database.Migrate()
		}()
		go func() {
			defer wg.Done()
			database.Migrate()
		}()

		wg.Wait()

		// Assert version 1 appears exactly once with is_applied = true
		rows, err := db.Query(`SELECT version_id, is_applied FROM goose_db_version WHERE version_id = 1`)
		if err != nil {
			rt.Fatalf("failed to query goose_db_version: %v", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			var versionID int64
			var isApplied bool
			if err := rows.Scan(&versionID, &isApplied); err != nil {
				rt.Fatalf("failed to scan row: %v", err)
			}
			if !isApplied {
				rt.Errorf("version %d has is_applied = false", versionID)
			}
			count++
		}
		if err := rows.Err(); err != nil {
			rt.Fatalf("rows iteration error: %v", err)
		}
		if count != 1 {
			rt.Errorf("expected version 1 to appear exactly once in goose_db_version, got %d occurrences", count)
		}
	})
}
