## 📜 Capuchin: A robust Todo application
A feature-rich full-stack todo list application with a Go (Golang) REST API backend and a React/Vite frontend using a professional-grade decoupled architecture.

## 🚀 Features Implemented

* **Backend (Go + Gin):** RESTful API with distinct layers (`handlers`, `services`, `database`, `middleware`) and robust error handling.
* **Authentication:** Secure Signup, Login, and Logout using short-lived JWT tokens with a database-backed token blacklisting mechanism.
* **Database (PostgreSQL):** Relational persistence mapped implicitly to user context to enforce cross-tenant data isolation.
* **Frontend (React + Vite):** Modern reactive UI with custom asynchronous Hooks (`useTodos`, `useAuth`) abstracting away native `fetch` requests.
* **Offline-friendly mode:** Supports an unauthenticated Guest mode backed tightly by `localStorage`.
* **Containerization:** Clean Docker Compose multi-stage orchestrations covering both isolated local development profiles and production scratch-image deployment.

## 📂 Project Structure

```text
capuchin/
├── backend/
│   ├── cmd/
│   │   ├── server/           # Entry point for the REST server
│   │   ├── migrate/          # Standalone binary runner for schema definitions
│   │   └── seed/             # Dev DB seed runner
│   ├── internal/
│   │   ├── config/           # Environment & Config map parsing
│   │   ├── database/         # PostgreSQL driver configuration & pooling limits 
│   │   ├── handlers/         # HTTP Route logic & payload validation
│   │   ├── middleware/       # Identity resolution & security guards
│   │   ├── models/           # Data structures
│   │   ├── routes/           # Mux mappings setup
│   │   └── services/         # Identity and persistence core logic workflows
│   ├── Dockerfile            # Multi-stage Backend Container 
│   ├── air.toml              # Hot Reload configs
│   ├── go.mod                # Go Dependencies
│   └── test.sh               # Integration / E2E endpoint bash test harness
├── frontend/
│   ├── src/
│   │   ├── components/       # Presentational layout components
│   │   ├── hooks/            # Primary React state workflows (`useAuth`, `useTodos`) 
│   │   ├── lib/              # Core native-fetch wrapper API logic
│   │   ├── pages/            # Page-level route views
│   │   ├── types/            # TypeScript definitions
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── Dockerfile            # Nginx + React Multi-stage Frontend Container
│   ├── vite.config.ts        # Vite bundling settings
│   └── package.json
├── compose.yml               # Lean Production Orchestration
├── compose-dev.yml           # Dev Mode (Air/Vite) overrides
└── Makefile                  # Command shortcuts
```

## 💻 Tech Stack
* **Backend:** Go (REST API, Clean Architecture)
* **Backend Framework:** Gin
* **Frontend:** React, TypeScript, Vite
* **Runtime Orchestration:** Docker, Make
* **Database:** PostgreSQL
* **Migrations:** Goose v3 (Inside Docker)

## 🛠️ How to Run

### Method 1: Docker (Recommended)
This approach encapsulates all dependencies securely via Docker Engine configurations.

#### For Development (Hot-Reloading)
Runs the Go backend natively through Air for hot-schema reload mappings, and the React frontend via Vite HMR.
```sh
make dev 
# OR
docker compose --env-file .env.example -f compose-dev.yml up --build
```
- **Frontend App**: `http://localhost:5173`
- **Backend API Base**: `http://localhost:8080`

#### For Production
Runs a lean production-ready sequence packaging the Go engine natively in a `scratch` container, and distributing the React codebase via `nginx`.
```bash
make prod
# OR
docker compose --env-file .env -f compose.yml up --build
```

### Method 2: Native via NPM script
Requires Go, Node.js, and Postgres installed natively on your machine!
Ensure your root `.env` accurately targets your native Postgres installation.
```bash
npm i
npx concurrently "cd ./backend/cmd/server && go run main.go" "npm run dev --prefix ./frontend"
```

---

## 🧠 Documentation & Key Concepts

For an in-depth dive into the structure, API contract, database schema, and best practices, please refer to our full documentation on the **[GitHub Wiki](https://github.com/the-monkeys/capuchin/wiki)**.

Key concepts utilized:
* **Go:** Structs, Slices, JSON Marshalling, Clean Architecture.
* **React:** Functional Components, Custom Hooks (`useTodos`, `useAuth`), fetch wrappers.
* **Testing:** `backend/test.sh` for E2E integration tests against API endpoints.
* **Docker:** Multi-stage builds, Scratch images, Docker Compose overrides.
* **General:** REST API Design, JWT Auth isolation, Postgres parameterization.

Long term plans:

folder todo
collaborators
real time update
organization
authentication
groups and access
sharelink
auth login
schedule with reminder
version control
mcp server
