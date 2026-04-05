# Capuchin Backend Architecture

This document outlines the architectural design and structural patterns used in the Capuchin Go backend.

## Overview

The backend is built using **Go** and the **Gin Web Framework**. It follows a variation of the **Clean Architecture** and the **Standard Go Project Layout**, ensuring separation of concerns, scalability, and maintainability.

The application interacts with a **PostgreSQL** database using the standard `database/sql` library and uses **JWT (JSON Web Tokens)** for stateless authentication.

## Directory Structure

The codebase is strictly divided into `cmd` for entry points and `internal` for private application code, preventing external imported usage of our core logic.

```text
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point. Wires dependencies.
├── internal/
│   ├── config/               # Environment loading and validation
│   ├── database/             # Global DB connection pool and schema init
│   ├── handlers/             # HTTP transport layer (Controllers)
│   ├── middleware/           # HTTP intercepts (Auth, Error Recovery)
│   ├── models/               # Domain data structures
│   ├── routes/               # Centralized route registration
│   └── services/             # Core business logic
└── ...
```

## Layered Architecture

The application handles requests through three primary layers:

1. **Routing Layer (`internal/routes`)**
   - Registers all endpoints to their corresponding handler functions.
   - Applies necessary middlewares (e.g., `AuthRequired`) to protected routes.

2. **Transport / Handler Layer (`internal/handlers`)**
   - Extracts and validates incoming HTTP requests (JSON body, Path params, Headers).
   - Calls the appropriate Service methods.
   - Formats the response (JSON) and returns appropriate HTTP status codes (200, 400, 404, 500).
   - **Rule:** Handlers contain *no business logic* or direct database queries.

3. **Service Layer (`internal/services`)**
   - Contains all the core business logic.
   - Enforces business rules (e.g., hashing passwords, verifying credentials, associating items).
   - Communicates directly with the data store (`internal/database`).
   - Returns business-level errors (e.g., `ErrUserExists`, `ErrTodoNotFound`) decoupled from HTTP transport.

## Dependency Injection

The application uses constructor injection to pass dependencies down the chain. This is primarily seen in the relationship between Handlers and Services:

```go
// main.go initializes components and wires them together
todoService := services.NewTodoService()
todoHandler := handlers.NewTodoHandler(todoService)
```

This decouples the handler from a strictly concrete service implementation, paving the way for easier unit testing via mocked services in the future.

## Database & Persistence

- **Connection Pool:** A centralized `sql.DB` connection pool (`database.DB`) is initialized at startup. It configures connection lifetimes, max open, and max idle connections to prevent resource exhaustion.
- **Schema Migrations:** Managed by [Goose v3](https://github.com/pressly/goose). `database.Migrate()` is called at startup and applies any pending SQL migrations in order. Migration files are embedded into the binary via `embed.FS`, making the binary fully self-contained with no external file dependencies at runtime. See [Backend Database Schema](backend_schema.md) for the full migration strategy and tooling rationale.
- **Relational Integrity:** Uses standard PostgreSQL relations (e.g., `todos.user_id REFERENCES users(id)`).
- **UUIDs:** Primary keys are decentralized using UUIDs.

## Authentication Flow

Authentication is stateless and managed via JWTs:

1. **Login:** A user logs in, the service verifies the hashed password via `bcrypt`, and generates an HS256 JWT containing the `user_id` and an expiration time.
2. **Authorization:** Protected routes use `middleware.AuthRequired()`, which intercepts requests, strictly validates the `Authorization` bearer token against the signing key, enforces the signing method, and extracts the `user_id` into the Gin context.
3. **Logout:** The application tracks revoked tokens using a database table `blacklisted_tokens`. When a user logs out, their specific token is inserted into this table. The auth middleware inherently rejects any blacklisted tokens. A background goroutine cleans up expired tokens hourly.
