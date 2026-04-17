# Implementation Plan: db-migration-ci-cleanup

## Overview

Pure refactoring/cleanup across five independent areas: extract a duplicated CI job into a reusable workflow, fix invalid port expressions in workflow files, strip commented-out deploy noise from `release.yml`, bump the Go version directive in `backend/migration/go.mod`, and consolidate `init.sql` to a single canonical path.

No runtime behaviour changes. Each task is self-contained and can be reviewed independently.

## Tasks

- [x] 1. Extract `regenerate-init-sql` into a reusable workflow
  - Create `.github/workflows/regenerate-init-sql.yml` triggered via `workflow_call` with an `environment` input
  - Copy the job body verbatim from the current inline job: checkout, install `postgresql-client`, `pg_dump` with `--schema-only --no-owner --no-privileges --exclude-table=goose_db_version`, prepend header, commit with `[skip ci]` and push
  - Use `${{ inputs.environment }}` in the commit message: `chore: regenerate init.sql after migration [${{ inputs.environment }}] [skip ci]`
  - Replace the inline `regenerate-init-sql` job in `release.yml` with a `uses:` call: `uses: ./.github/workflows/regenerate-init-sql.yml` with `environment: production` and `secrets: inherit`
  - Replace the inline `regenerate-init-sql` job in `migrate-manual.yml` with a `uses:` call: `uses: ./.github/workflows/regenerate-init-sql.yml` with `environment: ${{ inputs.environment }}`, `if: ${{ inputs.regenerate_init_sql }}`, and `secrets: inherit`
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [x] 2. Fix invalid port expressions in workflow files
  - In `release.yml` `migrate` job `env:` block, replace `${{ secrets.POSTGRES_PORT != '' && secrets.POSTGRES_PORT || '5432' }}` with `${{ secrets.POSTGRES_PORT }}`
  - In `migrate-manual.yml` `migrate` job `env:` block, apply the same replacement
  - Verify no workflow file contains the string `secrets.POSTGRES_PORT != ''` or the literal `5432`
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 3. Clean up commented-out deploy steps in `release.yml`
  - Remove the `Build backend Docker image` step from the `deploy` job
  - Remove the commented-out registry push example block and its associated inline comments
  - Remove the commented-out SSH deploy example block and its associated inline comments
  - Retain the `Deployment placeholder` step (`run: echo "Deploy capuchin-backend:..."`)
  - Retain the single explanatory comment directing developers to replace the placeholder with their actual deployment mechanism
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 4. Bump Go version directive in `backend/migration/go.mod`
  - Change `go 1.25.0` to `go 1.25.5` in `backend/migration/go.mod`
  - Run `go mod tidy` in `backend/migration/` and commit any resulting `go.sum` changes
  - _Requirements: 4.1, 4.3_

- [x] 5. Verify test-only dependencies in `backend/migration/go.mod`
  - Confirm that after `go mod tidy` (run in task 4), `github.com/google/uuid` and `golang.org/x/crypto` are listed correctly — if `go mod tidy` keeps them as direct dependencies (because test-file imports count), no further change is needed; the module graph is already honest
  - If `go mod tidy` moves them to `// indirect`, commit the updated `go.mod`
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 6. Consolidate `init.sql` to a single canonical file
  - Update `backend/migration/init_sql_test.go`: change `os.ReadFile("db/init.sql")` to `os.ReadFile("../../db/init.sql")`
  - Delete `backend/migration/db/init.sql`
  - Confirm `backend/migration/db/embed.go` is unchanged (it only embeds `migrations/*.sql` via glob and never matched `init.sql`)
  - _Requirements: 6.1, 6.2, 6.3, 6.5, 6.7_

- [x] 7. Checkpoint — verify build, vet, and unit tests pass
  - Run `go build ./...` and `go vet ./...` in `backend/migration/` — must produce no errors
  - Run `go test ./...` in `backend/` — must pass with no failures
  - Ensure all tests pass, ask the user if questions arise.
  - _Requirements: 7.1, 7.2_

## Notes

- No task is optional
- Tasks 1–3 are CI YAML changes; tasks 4–5 are `go.mod` changes; task 6 is a file deletion + test path fix
- `TestInitSQLMatchesMigrationEndState` (Requirement 7.3) requires Docker via testcontainers and runs only on the weekly CI schedule — it is not part of the local `go test ./...` run
- No property-based tests apply: this spec contains no pure functions with a meaningful input space
