ifneq (, $(shell command -v podman 2> /dev/null))
	CONTAINER_RUNTIME := podman
endif

ifneq (, $(shell command -v docker 2> /dev/null))
	CONTAINER_RUNTIME := docker
endif

# Docker Dev Mode (Hot Reload)
dev:
	$(CONTAINER_RUNTIME) compose --env-file .env.example -f compose-dev.yml up --build -d

dev-logs: 
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml logs 

dev-down:
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml down
clean:
	$(CONTAINER_RUNTIME) compose -f compose-dev.yml down --volumes --remove-orphans --rmi all


prod:
	$(CONTAINER_RUNTIME) compose --env-file .env -f compose.yml up

logs:
	$(CONTAINER_RUNTIME) compose -f compose.yml logs -f

down:
	$(CONTAINER_RUNTIME) compose -f compose.yml down


frontend:
	cd frontend && npm run dev

backend:
	cd backend && air

seed:
	cd backend && go run ./cmd/seed

.PHONY:  dev dev-logs dev-down prod logs down seed

