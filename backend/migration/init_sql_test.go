package migration_test

// TestInitSQLMatchesMigrationEndState verifies that migration/db/init.sql
// produces a schema structurally identical to the one produced by running all
// goose migrations. This catches cases where init.sql was not regenerated after
// a new migration was added, or was accidentally hand-edited.
//
// Feature: migration-module-separation, Property 9: init.sql schema matches goose migration end state.

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// newNamedTestDB spins up a fresh postgres testcontainer with the given DB name.
func newNamedTestDB(t *testing.T, dbName string) *sql.DB {
	t.Helper()
	ctx := context.Background()
	pgc, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container (%s): %v", dbName, err)
	}
	t.Cleanup(func() { _ = pgc.Terminate(ctx) })

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string (%s): %v", dbName, err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open db (%s): %v", dbName, err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// tableSchema holds the structural description of a single table.
type tableSchema struct {
	columns     []columnDef
	constraints []constraintDef
}

type columnDef struct {
	name       string
	dataType   string
	isNullable string
}

type constraintDef struct {
	name           string
	constraintType string
}

// dumpSchema queries information_schema for all user tables, their columns,
// and their constraints, returning a normalised map keyed by table name.
func dumpSchema(t *testing.T, db *sql.DB) map[string]tableSchema {
	t.Helper()

	// Fetch tables.
	tableRows, err := db.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_type = 'BASE TABLE'
		  AND table_name NOT LIKE 'goose_%'
		ORDER BY table_name`)
	if err != nil {
		t.Fatalf("dumpSchema: query tables: %v", err)
	}
	defer tableRows.Close()

	schema := make(map[string]tableSchema)
	for tableRows.Next() {
		var name string
		if err := tableRows.Scan(&name); err != nil {
			t.Fatalf("dumpSchema: scan table name: %v", err)
		}
		schema[name] = tableSchema{}
	}
	if err := tableRows.Err(); err != nil {
		t.Fatalf("dumpSchema: table rows error: %v", err)
	}

	// Fetch columns for each table.
	for tableName, ts := range schema {
		colRows, err := db.Query(`
			SELECT column_name, data_type, is_nullable
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = $1
			ORDER BY ordinal_position`, tableName)
		if err != nil {
			t.Fatalf("dumpSchema: query columns for %s: %v", tableName, err)
		}
		var cols []columnDef
		for colRows.Next() {
			var c columnDef
			if err := colRows.Scan(&c.name, &c.dataType, &c.isNullable); err != nil {
				colRows.Close()
				t.Fatalf("dumpSchema: scan column: %v", err)
			}
			cols = append(cols, c)
		}
		colRows.Close()
		if err := colRows.Err(); err != nil {
			t.Fatalf("dumpSchema: column rows error: %v", err)
		}
		ts.columns = cols

		// Fetch constraints.
		conRows, err := db.Query(`
			SELECT constraint_name, constraint_type
			FROM information_schema.table_constraints
			WHERE table_schema = 'public' AND table_name = $1
			ORDER BY constraint_name`, tableName)
		if err != nil {
			t.Fatalf("dumpSchema: query constraints for %s: %v", tableName, err)
		}
		var cons []constraintDef
		for conRows.Next() {
			var c constraintDef
			if err := conRows.Scan(&c.name, &c.constraintType); err != nil {
				conRows.Close()
				t.Fatalf("dumpSchema: scan constraint: %v", err)
			}
			cons = append(cons, c)
		}
		conRows.Close()
		if err := conRows.Err(); err != nil {
			t.Fatalf("dumpSchema: constraint rows error: %v", err)
		}
		ts.constraints = cons
		schema[tableName] = ts
	}

	return schema
}

// schemaKey produces a deterministic string representation of a schema map
// for easy diffing in test output.
func schemaKey(schema map[string]tableSchema) string {
	tables := make([]string, 0, len(schema))
	for t := range schema {
		tables = append(tables, t)
	}
	sort.Strings(tables)

	var sb strings.Builder
	for _, t := range tables {
		ts := schema[t]
		fmt.Fprintf(&sb, "TABLE %s\n", t)
		for _, c := range ts.columns {
			fmt.Fprintf(&sb, "  COL %s %s nullable=%s\n", c.name, c.dataType, c.isNullable)
		}
		cons := make([]string, len(ts.constraints))
		for i, c := range ts.constraints {
			cons[i] = fmt.Sprintf("%s:%s", c.constraintType, c.name)
		}
		sort.Strings(cons)
		for _, c := range cons {
			fmt.Fprintf(&sb, "  CON %s\n", c)
		}
	}
	return sb.String()
}

func TestInitSQLMatchesMigrationEndState(t *testing.T) {
	// Feature: migration-module-separation, Property 9: init.sql schema matches goose migration end state

	// Read init.sql from disk (relative to the migration/ module root).
	initSQL, err := os.ReadFile("db/init.sql")
	if err != nil {
		t.Fatalf("failed to read db/init.sql: %v", err)
	}

	// DB 1: bootstrapped via init.sql
	dbInit := newNamedTestDB(t, "testinit")
	if _, err := dbInit.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to apply init.sql: %v", err)
	}

	// DB 2: bootstrapped via goose migrations
	dbGoose := newNamedTestDB(t, "testgoose")
	migrateDB(t, dbGoose)

	// Compare schemas.
	schemaInit := dumpSchema(t, dbInit)
	schemaGoose := dumpSchema(t, dbGoose)

	keyInit := schemaKey(schemaInit)
	keyGoose := schemaKey(schemaGoose)

	if keyInit != keyGoose {
		t.Errorf("init.sql schema does not match goose migration end state\n\n--- init.sql ---\n%s\n--- goose ---\n%s", keyInit, keyGoose)
	}
}
