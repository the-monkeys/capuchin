ifneq (, $(shell command -v podman 2> /dev/null))
	CONTAINER_RUNTIME := podman
endif

ifneq (, $(shell command -v docker 2> /dev/null))
	CONTAINER_RUNTIME := docker
endif

.DEFAULT_GOAL := help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Dev (hot reload via Docker) ───────────────────────────────────────────────

dev: ## Start all services in dev mode (hot reload)
	$(CONTAINER_RUNTIME) compose --env-file .env.example -f compose-dev.yml up --build -d

dev-logs: ## Tail dev logs
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml logs -f

dev-down: ## Stop dev services
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml down

clean: ## Stop dev services and remove volumes, images, orphans
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml down --volumes --remove-orphans --rmi all

# ── Prod ──────────────────────────────────────────────────────────────────────

prod: ## Start all services in prod mode (detached)
	$(CONTAINER_RUNTIME) compose --env-file .env -f compose.yml up -d

logs: ## Tail prod logs
	$(CONTAINER_RUNTIME) compose -f compose.yml logs -f

down: ## Stop prod services
	$(CONTAINER_RUNTIME) compose -f compose.yml down

# ── Local dev (outside Docker) ────────────────────────────────────────────────

frontend: ## Start frontend dev server
	cd frontend && npm run dev

backend: ## Start backend with hot reload (requires air: go install github.com/air-verse/air@v1.61.7)
	cd backend && air

# ── Database ──────────────────────────────────────────────────────────────────

migrate: ## Run migrations against localhost DB (reads .env for credentials)
	@set -a && . ./.env.example && set +a && export POSTGRES_HOST=localhost && cd backend/migration && go run ./cmd/migrate up

migrate-down: ## Roll back the last migration against localhost DB
	@set -a && . ./.env.example && set +a && export POSTGRES_HOST=localhost && cd backend/migration && go run ./cmd/migrate down

seed: ## Seed dev database with sample data (reads .env.example for credentials)
	@set -a && . ./.env.example && set +a && export POSTGRES_HOST=localhost && cd backend && go run ./cmd/seed

migrate-build: ## Build migration Docker image
	docker build -f backend/migration/Dockerfile -t capuchin-migration ./backend

.PHONY: help dev dev-logs dev-down clean prod logs down frontend backend migrate migrate-down seed migrate-build
