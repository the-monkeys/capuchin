ifneq (, $(shell command -v podman 2> /dev/null))
	CONTAINER_RUNTIME := podman
endif

ifneq (, $(shell command -v docker 2> /dev/null))
	CONTAINER_RUNTIME := docker
endif

# Docker Dev Mode (Hot Reload)
dev:
	$(CONTAINER_RUNTIME) compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d

logs:
	$(CONTAINER_RUNTIME) compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

dev-down:
	$(CONTAINER_RUNTIME) compose -f docker-compose.yml -f docker-compose.dev.yml down

# Docker Production Mode
prod:
	$(CONTAINER_RUNTIME) compose up --build

# Stop Containers
down:
	$(CONTAINER_RUNTIME) compose down

frontend:
	cd frontend && npm run dev

backend:
	cd backend && air

.PHONY: frontend backend dev prod down

