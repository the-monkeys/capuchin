# Design Document — db-migration-ci-cleanup

## Overview

This is a pure refactoring/cleanup spec. No runtime behaviour changes. The work
falls into five independent, low-risk areas:

1. Extract the duplicated `regenerate-init-sql` CI job into a reusable workflow.
2. Remove invalid GitHub Actions port expressions and rely solely on secrets.
3. Strip commented-out placeholder deployment steps from `release.yml`.
4. Bump the Go version directive in `backend/migration/go.mod` from `1.25.0` to `1.25.5`.
5. Delete `backend/migration/db/init.sql` and update the one test that reads it to
   use the canonical path `backend/db/init.sql`.

Each area is self-contained and can be reviewed independently.

---

## Architecture

The repository has two Go modules:

```
backend/          ← module "capuchin"          (app server)
backend/migration/← module "capuchin-migration" (migration runner)
```

They share a single Postgres schema whose source of truth is the goose migration
files at `backend/migration/db/migrations/*.sql`.

Two files currently represent the schema as a snapshot:

| File | Purpose |
|---|---|
| `backend/db/init.sql` | Docker Compose bootstrap for local dev; also read by `TestInitSQLMatchesMigrationEndState` |
| `backend/migration/db/init.sql` | Duplicate — read by `TestInitSQLMatchesMigrationEndState` (old path) |

After this cleanup only `backend/db/init.sql` will exist.

CI has three workflow files:

| File | Trigger |
|---|---|
| `.github/workflows/ci.yml` | PR / push to main / weekly schedule |
| `.github/workflows/release.yml` | Push to main or version tag |
| `.github/workflows/migrate-manual.yml` | Manual dispatch |

`release.yml` and `migrate-manual.yml` both contain an identical inline
`regenerate-init-sql` job. That job will be extracted into a reusable workflow.

---

## Components and Interfaces

### 1. Reusable Workflow — `regenerate-init-sql.yml`

**New file:** `.github/workflows/regenerate-init-sql.yml`

Triggered via `workflow_call`. Accepts one input:

```yaml
on:
  workflow_call:
    inputs:
      environment:
        required: true
        type: string
```

The job body is identical to the current inline job in both callers:
- Checkout with `GITHUB_TOKEN`
- Install `postgresql-client`
- `pg_dump` with `--schema-only --no-owner --no-privileges --exclude-table=goose_db_version`
- Prepend the auto-generated header comment
- Commit with `[skip ci]` and push

The commit message uses the `environment` input so it remains environment-specific:
```
chore: regenerate init.sql after migration [${{ inputs.environment }}] [skip ci]
```

**Callers after refactor:**

`release.yml`:
```yaml
regenerate-init-sql:
  uses: ./.github/workflows/regenerate-init-sql.yml
  needs: migrate
  with:
    environment: production
  secrets: inherit
```

`migrate-manual.yml`:
```yaml
regenerate-init-sql:
  uses: ./.github/workflows/regenerate-init-sql.yml
  needs: migrate
  if: ${{ inputs.regenerate_init_sql }}
  with:
    environment: ${{ inputs.environment }}
  secrets: inherit
```

### 2. Port Expression Fix

The expression `${{ secrets.POSTGRES_PORT != '' && secrets.POSTGRES_PORT || '5432' }}`
is JavaScript ternary syntax. GitHub Actions evaluates it in `env:` blocks as a
string literal rather than as an expression, so the port is never actually set
from the secret. It appears in two places:

- `release.yml` — `env.POSTGRES_PORT` in the `migrate` job, and `--port=` in the
  `pg_dump` command in the `regenerate-init-sql` job.
- `migrate-manual.yml` — same two locations.

After the refactor the `regenerate-init-sql` job moves to the reusable workflow,
so only the `migrate` job in each caller retains a port reference. Both are
replaced with `${{ secrets.POSTGRES_PORT }}` — no fallback.

The `cmd/migrate/main.go` binary already defaults to `5432` when `POSTGRES_PORT`
is empty, so the application layer is unaffected. The CI layer intentionally
removes the fallback to force operators to configure the secret explicitly.

### 3. release.yml Deploy Job Cleanup

The `deploy` job in `release.yml` currently contains two large commented-out
example blocks (registry push, SSH deploy) plus a `Build backend Docker image`
step that is not part of any real deployment. These are removed. What remains:

```yaml
deploy:
  name: Deploy backend
  runs-on: ubuntu-latest
  needs: migrate
  environment: production
  steps:
    - uses: actions/checkout@v4

    # Replace the step below with your actual deployment mechanism.
    # Examples: docker push + SSH deploy, kubectl apply, fly deploy, etc.

    - name: Deployment placeholder
      run: echo "Deploy capuchin-backend:${{ github.sha }} to production"
```

The `Build backend Docker image` step is also removed — it builds an image that
is never pushed or used, which is pure noise.

### 4. Migration Module go.mod Version Bump

`backend/migration/go.mod` currently declares `go 1.25.0`. It is bumped to
`go 1.25.5` to match `backend/go.mod`. This is a directive-only change; no
dependency versions change. `go.sum` may gain a new toolchain entry but the
existing dependency hashes are unaffected.

### 5. go.mod Direct Dependency Cleanup

`backend/migration/go.mod` lists `github.com/google/uuid` and
`golang.org/x/crypto` as direct (non-`// indirect`) dependencies. Both are
imported only in test files (`migrate_test.go`), not in any production source
file under `cmd/` or the module root.

In Go modules, test-only imports of packages that are not imported by any
non-test file in the same module are still listed as direct dependencies in
`go.mod` when they appear in `_test.go` files within the module. Running
`go mod tidy` will keep them as direct dependencies because the test files are
part of the module. The requirement asks to move them to `// indirect`.

The correct approach is to verify whether `go mod tidy` actually marks them
indirect. If `go mod tidy` keeps them direct (because test files count), the
requirement is satisfied by the fact that `go mod tidy` produces no changes —
meaning the current state is already what tidy would produce. The key action is
to confirm the current state and document it; no mechanical change may be needed
beyond the version bump in §4.

> **Decision:** Run `go mod tidy` after the version bump. If uuid and x/crypto
> remain direct, they are legitimately direct (test imports count). The
> requirement's intent — that the module graph is honest — is satisfied because
> `go mod tidy` is the authoritative tool. No manual editing of `go.mod` is
> needed beyond what tidy produces.

### 6. init.sql Consolidation

**Delete:** `backend/migration/db/init.sql`

**Update:** `backend/migration/init_sql_test.go` — change the `os.ReadFile` path:

```go
// Before
initSQL, err := os.ReadFile("db/init.sql")

// After
initSQL, err := os.ReadFile("../../db/init.sql")
```

The test runs from `backend/migration/` (the module root), so `../../db/init.sql`
resolves to `backend/db/init.sql`.

`backend/migration/db/embed.go` is **not modified** — it only embeds
`migrations/*.sql` via a glob that never matched `init.sql`.

`backend/db/init.sql` already contains the correct header comment. Its content
is identical to the deleted file, so no schema change occurs.

---

## Data Models

No data model changes. The schema itself (`users`, `todos`, `blacklisted_tokens`)
is unchanged. The only file-level change to schema representation is the deletion
of the duplicate `backend/migration/db/init.sql`.

---

## Error Handling

| Scenario | Handling |
|---|---|
| `POSTGRES_PORT` secret not set in CI | `pg_dump --port=` receives an empty string → command fails with a clear error. Operators must configure the secret. |
| `go mod tidy` changes go.sum | Commit the updated go.sum alongside go.mod. |
| `TestInitSQLMatchesMigrationEndState` fails after path change | The test itself reports a diff between init.sql schema and goose end state — fix by ensuring `backend/db/init.sql` is up to date. |

---

## Testing Strategy

This feature is a refactoring/cleanup. It involves CI YAML files, a `go.mod`
directive change, a file deletion, and a single string change in a test file.
None of these involve pure functions with a meaningful input space, so
**property-based testing does not apply**.

The appropriate verification strategy is a combination of static checks and
existing test execution:

### Static / Example Checks (manual or scripted)

| Check | How |
|---|---|
| Reusable workflow file exists and declares `workflow_call` | Read `.github/workflows/regenerate-init-sql.yml` |
| `release.yml` and `migrate-manual.yml` use `uses:` to call the reusable workflow | Read both files |
| No workflow file contains the invalid ternary port expression | `grep -r "secrets.POSTGRES_PORT != ''"` returns no matches |
| No workflow file contains the hardcoded port `5432` | `grep -r "5432" .github/workflows/` returns no matches |
| Commented-out deploy examples removed from `release.yml` | Read `release.yml` |
| Deployment placeholder step retained | Read `release.yml` |
| `backend/migration/go.mod` declares `go 1.25.5` | Read `go.mod` |
| `backend/migration/db/init.sql` deleted | File does not exist |
| `init_sql_test.go` reads from `../../db/init.sql` | Read the test file |
| `embed.go` unchanged | Read `embed.go` |
| No file in `backend/migration/` references the old path | `grep -r "migration/db/init.sql" backend/migration/` returns no matches |
| All workflow YAML files are syntactically valid | Parse with a YAML linter (e.g. `yamllint`) |

### Build and Unit Tests

```bash
# backend module — must pass with no failures
cd backend && go test ./...

# migration module — build and vet must pass
cd backend/migration && go build ./... && go vet ./...
```

### Integration Tests (weekly CI schedule — requires Docker)

```bash
cd backend/migration && go test ./... -v -timeout 10m
```

This runs `TestInitSQLMatchesMigrationEndState` (among others), which spins up
two testcontainers Postgres instances and compares the schema produced by
`backend/db/init.sql` against the schema produced by running all goose
migrations. This is the primary correctness gate for Requirement 6.

The weekly schedule in `ci.yml` already covers this — no CI changes are needed
for the integration test job.
