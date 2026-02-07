## 📜 Capuchin: A basic Todo app
A basic full-stack todo list application with a Go (Golang) REST API backend and a React frontend with a professional-grade storage architecture.

## 🚀 Features Implemented

* **Backend (Go + Gin):** RESTful API with CRUD operations.
* **Frontend (React + Vite):** Modern reactive UI with Hooks (useState, useEffect).
* **Styling (Tailwind CSS):** Dark-mode interface with optimistic UI.
* **Persistence:** File-based JSON storage.
* **Architecture:** Refactored into "Standard Go Layout" (cmd, internal).
* **Containerization:** Docker & Docker Compose for Dev/Prod.


## 📂 Project Structure

```
capuchin/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Entry point
│   ├── db/
│   │   └── db.json               # Database
│   ├── internal/
│   │   ├── api/
│   │   ├── models/
│   │   │   └── todo.go           # Data structures
│   │   └── store/
│   │       └── file.go           # File I/O logic
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
├── docker-compose.yml            # Prod Orchestration
├── docker-compose.dev.yml        # Dev Mode Overrides
└── makefile                      # Command shortcuts
└── package.json                     

```

## 💻 Tech Stack
* **Backend:** Go (REST API)
* **Frontend:** React, TypeScript
* **Containerize:** Docker
* **Database:** File System storage
* **Backend Framework:** Gin
* **Frontend:** React, TypeScript
* **Containerize:** Docker
* **Database:** File System storage

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

🔴 Hot reload not working yet (for backend)

```bash
make dev
# OR
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```
- **Frontend**: http://localhost:5173
- **Health Check**: http://localhost:8080/health
- **Backend API**: http://localhost:8080/todos

#### <u > Production Mode</u >
Runs a lean, production-ready build (`scratch` image for Go, `nginx` for React).

```bash
make prod
# OR
docker compose up --build
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

* **Go:** Structs, Slices, JSON Marshalling, Modules, Package Exporting.
* **React:** Functional Components, Hooks, API Integration (fetch, async/await), Controlled Inputs.
* **Docker:** Multi-stage builds, Scratch images, Docker Compose overrides.
* **General:** REST API Design, CORS, JSON Persistence, Refactoring,TypeScript(for styling), axios (for API calls)


# **Plans for v1:**
authentication
postgres db (with ORM?)
units and sub-units (folders)

# **Long term plans:**
mcp server
collaborators
real time update
organization
groups and access
sharelink
auth login
schedule with reminder
version control
mcp server

