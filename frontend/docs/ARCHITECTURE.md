# Frontend Architecture Documentation

## Overview

This document describes the architecture of the Capuchin frontend application - a React-based todo management app with authentication support.

## Tech Stack

| Category | Technology |
|----------|------------|
| Framework | React 19 |
| Language | TypeScript 5.9 |
| Build Tool | Vite (rolldown-vite) |
| Routing | React Router v7 |
| Styling | Tailwind CSS v4 |
| HTTP Client | Native Fetch API |
| State Management | React Hooks (useState, useEffect, useCallback) |

---

## Project Structure

```
src/
├── assets/              # Static assets (images, fonts, icons)
│   └── capuchin.svg
├── components/         # Reusable React components
│   ├── layout/        # Layout components
│   │   └── Navbar.tsx
│   ├── todos/        # Todo-specific components
│   │   ├── TodoFilters.tsx
│   │   ├── TodoInput.tsx
│   │   └── TodoItem.tsx
│   └── ui/           # Generic UI primitives
│       ├── Alert.tsx
│       ├── Button.tsx
│       ├── EmptyState.tsx
│       └── Input.tsx
├── config/            # Configuration files
│   └── env.ts         # Environment variables
├── hooks/             # Custom React hooks
│   ├── useAuth.ts     # Authentication logic
│   └── useTodos.ts    # Todo state management
├── lib/               # Library utilities
│   └── api.ts         # API client functions
├── pages/             # Route page components
│   ├── Login.tsx
│   ├── Signup.tsx
│   └── Todos.tsx
├── types/             # TypeScript type definitions
│   └── index.ts
├── App.css            # Global styles
├── App.tsx            # Root component with routing
└── main.tsx           # Application entry point
```

---

## Component Hierarchy

```
App
├── Routes
│   ├── Todos (/)
│   │   ├── Navbar
│   │   ├── TodoInput
│   │   ├── TodoFilters
│   │   └── TodoItem[]
│   ├── Login (/login)
│   └── Signup (/signup)
```

---

## Data Flow

### Authentication Flow

```
User Login
    │
    ▼
authApi.login(email, password)
    │
    ▼ (on success)
saveToken(token) ──▶ localStorage.setItem("token")
    │
    ▼
setToken(state)
    │
    ▼
isAuthed = true
```

### Dual-Mode System

The application supports two modes:

#### 1. Guest Mode (Unauthenticated)
- Todos stored in localStorage under key `capuchin_guest_todos`
- All operations are local-only
- No JWT token required

#### 2. Authenticated Mode
- Todos stored on server
- JWT token required for all API calls
- Token stored in localStorage under key `token`

```
┌─────────────────────────────────────┐
│           useTodos Hook             │
├─────────────────────────────────────┤
│  isAuthed = false?                  │
│       │                             │
│       ├── YES ──▶ Guest Mode        │
│       │         (localStorage)      │
│       │                             │
│       └── NO ──▶ Authenticated     │
│                  (API calls)        │
└─────────────────────────────────────┘
```

### Todo Operations Flow

```
User Action (add/toggle/delete/update)
         │
         ▼
    useTodos hook
         │
         ├──▶ Optimistic UI Update (immediate)
         │
         └──▶ API Call (or localStorage)
                  │
                  ├──▶ Success  → Keep changes
                  │
                  └──▶ Failure  → Rollback to snapshot
```

---

## API Integration Pattern

### Centralized API Client

All API calls are centralized in `src/lib/api.ts`:

```typescript
// API modules export typed methods
export const authApi = {
  login: async (email: string, password: string): Promise<{ token: string }> => { ... },
  signup: async (email: string, password: string): Promise<void> => { ... },
}

export const todosApi = {
  getAll: async (token: string): Promise<Todo[]> => { ... },
  create: async (token: string, item: string): Promise<Todo> => { ... },
  toggle: async (token: string, id: string): Promise<void> => { ... },
  update: async (token: string, id: string, item: string): Promise<Todo> => { ... },
  delete: async (token: string, id: string): Promise<void> => { ... },
}
```

### Authentication Headers

```typescript
const authHeaders = (token: string): HeadersInit => ({
  "Content-Type": "application/json",
  Authorization: `Bearer ${token}`,
})
```

### Data Normalization

API responses are normalized to match frontend types:

```typescript
const normalise = (raw: Record<string, unknown>): Todo => ({
  id: String(raw.id ?? raw.ID),
  item: String(raw.item ?? raw.Item ?? ""),
  completed: Boolean(raw.completed ?? raw.Completed ?? false),
})
```

---

## State Management

### Custom Hooks Pattern

State is managed through custom hooks following React best practices:

#### `useAuth` Hook
- Manages authentication state
- Token storage/retrieval from localStorage
- JWT token parsing (extracts email from payload)
- Logout functionality

#### `useTodos` Hook
- Manages todo list state
- Handles both guest and authenticated modes
- Provides CRUD operations
- Implements optimistic updates with rollback
- Handles filtering (all/active/done)

---

## Routing

Routes are defined in `App.tsx` using React Router:

```typescript
<Routes>
  <Route path="/" element={<Todos />} />
  <Route path="/todos" element={<Todos />} />
  <Route path="/signup" element={<Signup />} />
  <Route path="/login" element={<Login />} />
</Routes>
```

---

## Configuration

### Environment Variables

Environment configuration is centralized in `src/config/env.ts`:

```typescript
export const env = {
  apiUrl: import.meta.env.VITE_API_URL,
} as const
```

### TypeScript Path Aliases

Configured in `tsconfig.app.json`:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  }
}
```

This enables clean imports like:
```typescript
import { Button } from "@/components/ui/Button"
import { useAuth } from "@/hooks/useAuth"
```

---

## Styling Strategy

### Tailwind CSS v4

The project uses Tailwind CSS v4 with:
- CSS-first configuration
- Custom utility classes for components
- Responsive design with breakpoints
- Consistent spacing and sizing

### Component Variants

UI components support variants and sizes:

```typescript
// Button component example
type Variant = "primary" | "secondary" | "ghost" | "danger"
type Size = "sm" | "md" | "lg"

interface ButtonProps {
  variant?: Variant
  size?: Size
  fullWidth?: boolean
  loading?: boolean
}
```

---

## Error Handling

### API Errors

- Thrown as descriptive errors with HTTP status context
- Caught and logged to console in hooks
- UI reflects error states through loading/empty states

### Optimistic Updates with Rollback

```typescript
const deleteTodo = useCallback(async (id: string) => {
  // Snapshot for potential rollback
  const snapshot = todos
  
  // Immediate UI update
  setTodos((prev) => prev.filter((t) => t.id !== id))
  
  try {
    await todosApi.delete(token!, id)
  } catch {
    // Rollback on failure
    setTodos(snapshot)
  }
}, [isAuthed, token, todos, updateGuest])
```

---

## Security Considerations

### Token Storage
- JWT tokens stored in localStorage
- Token payload decoded client-side for user info
- Tokens cleared on logout

### API Security
- Bearer token authentication
- Content-Type: application/json for all requests

---

## Build & Development

### Available Scripts

| Script | Command | Description |
|--------|---------|-------------|
| dev | `npm run dev` | Start development server |
| build | `npm run build` | Type-check and build for production |
| lint | `npm run lint` | Run ESLint |
| preview | `npm run preview` | Preview production build |

---