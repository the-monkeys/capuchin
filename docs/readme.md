## 📜 Capuchin: A basic Todo app
A basic full-stack todo list application with a Go (Golang) REST API backend and a React frontend with a professional-grade storage architecture.

## 🚀 Features Implemented

* **Backend (Go + Gin):** RESTful API with distinct layers (Handlers, Services, DB) and robust error handling.
* **Authentication:** Secure Signup, Login, and Logout using JWT tokens.
* **Database (PostgreSQL):** Relational persistence using `database/sql` with schema migrations managed by Goose v3 on startup.
* **Frontend (React + Vite):** Modern reactive UI with Hooks (useState, useEffect).
* **Styling (Tailwind CSS):** Dark-mode interface with optimistic UI.
* **Architecture:** Clean architecture enforcing separation of concerns in 'internal'.
* **Containerization:** Docker & Docker Compose for Dev/Prod.

## 📂 Project Structure

```
capuchin/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Entry point
│   ├── internal/
│   │   ├── config/               # Environment & Config setup
│   │   ├── database/             # PostgreSQL connection & init
│   │   ├── handlers/             # HTTP Route handlers
│   │   ├── middleware/           # Auth & Error middleware
│   │   ├── models/               # Data structures
│   │   ├── routes/               # API route definitions
│   │   └── services/             # Core business logic
│   ├── Dockerfile                # Backend Container
│   ├── air.toml                  # Hot Reload Config
│   ├── go.mod                    # Dependencies
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   ├── App.css
│   │   └── main.tsx
│   ├── Dockerfile                # Frontend Container
│   ├── vite.config.ts            # Build Config
│   └── package.json
├── compose.yml                   # Prod Orchestration
├── compose-dev.yml               # Dev Mode Overrides
└── Makefile                      # Command shortcuts
└── package.json                     

```

## 💻 Tech Stack
## 💻 Tech Stack
* **Backend:** Go (REST API, Clean Architecture)
* **Backend Framework:** Gin
* **Frontend:** React, TypeScript
* **Containerize:** Docker
* **Database:** PostgreSQL
* **Migrations:** Goose v3

## 🛠️ How to Run

### Method 1:   In separate terminals



#### <b> Backend: </b>

Open Terminal 1
``` Bash
cd backend
go run cmd/server/main.go
```
`Server runs on localhost:8080`

#### Frontend:

Open Terminal 2
``` Bash
cd frontend
npm run dev
```
`Client opens at localhost:5173`


---

### Method 2:   Using npm Script (In project home directory)

Install npm packages
``` Bash
npm i
```
Run npx script

``` Bash
npx concurrently "cd ./backend/cmd/server && go run main.go" "npm run dev --prefix ./frontend"
```

- **Frontend**: http://localhost:5173
- **Health Check**: http://localhost:8080/health
- **Backend API**: http://localhost:8080/todos


---


### Method 3: Docker (In project home directory)

We support two modes: **Development** (Hot-Reload) and **Production** (Lean Static Builds).

#### <u >Development Mode</u >
Runs the backend with `Air` (Go hot-reload) and Frontend with `Vite` (HMR). Changes to code are reflected instantly.

```bash
make dev
# OR
docker compose --env-file .env.example -f compose-dev.yml up --build
```
- **Frontend**: http://localhost:5173
- **Health Check**: http://localhost:8080/health
- **Backend API**: http://localhost:8080/todos

#### <u > Production Mode</u >
Runs a lean, production-ready build (`scratch` image for Go, `nginx` for React).

```bash
make prod
# OR
docker compose --env-file .env -f compose.yml up --build
```
- **App**: http://localhost
- **Health Check**: http://localhost:8080/health
- **Backend API**: http://localhost:8080/todos

#### Stop Containers
```bash
make down
# OR
#in active terminal
ctrl+c or cmd+c 
```

---



    
## 🧠 Key Concepts Implemented (can be seen in comments)

For an in-depth dive into the structure and patterns, please refer to our dedicated documentation:
- [Backend Architecture Reference](backend_architecture.md)
- [Backend Best Practices](backend_best_practices.md)
- [Backend API Contract](backend_api.md)
- [Backend Database Schema](backend_schema.md)

* **Go:** Structs, Slices, JSON Marshalling, Modules, Package Exporting, Clean Architecture.
* **React:** Functional Components, Hooks, API Integration (fetch, async/await), Controlled Inputs.
* **Testing:** Included a robust `backend/verify_backend.sh` shell script to instantly orchestrate E2E integration tests against all API endpoints.
* **Docker:** Multi-stage builds, Scratch images, Docker Compose overrides.
* **General:** REST API Design, CORS, JSON Persistence, Refactoring,TypeScript(for styling), axios (for API calls)


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
