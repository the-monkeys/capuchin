## Capuchin: A basic Todo app
A basic full-stack todo list application with a Go (Golang) REST API backend and a React frontend with a professional-grade storage architecture.

## 🚀 Features Implemented

Backend (Go + Gin): RESTful API with CRUD operations.
Frontend (React + Vite): Modern reactive UI with Hooks (useState, useEffect).
Styling (Tailwind CSS): Dark-mode interface with optimistic UI.
Persistence: File-based JSON storage.
Architecture: Refactored into "Standard Go Layout" (cmd, internal).


## 📂 Final Project Structure

```capuchin/
├── backend/
│   ├── cmd/server/main.go       # Entry point
│   ├── internal/models/todo.go  # Data structures
│   ├── internal/store/file.go   # File I/O logic
│   ├── db/db.json               # Database
│   ├── go.mod                   # Dependencies
│   └── go.sum
└── frontend/
    ├── src/App.tsx              # React UI & Logic
    ├── src/index.css            # Tailwind Imports
    ├── vite.config.ts           # Build Config
    └── package.json
```

## Tech stack
* **Backend:** Go (REST API).
* Backend Framework: Gin
* **Frontend:** React, TypeScript
* **Database:** File System storage

## 🛠️ How to Run

### Backend:

Open Terminal 1
``` Bash
cd backend
go run cmd/server/main.go
```
Server runs on localhost:8080

### Frontend:

Open Terminal 2
``` Bash
cd frontend
npm run dev
```
Browser opens at localhost:5173


    
## 🧠 Key Concepts Implemented (can be seen in comments)

* **Go:** Structs, Slices, JSON Marshalling, Modules, Package Exporting.
* **React:** Functional Components, Hooks, API Integration (fetch, async/await), Controlled Inputs.
* **General:** REST API Design, CORS, JSON Persistence, Refactoring,TypeScript(for styling), axios (for API calls)