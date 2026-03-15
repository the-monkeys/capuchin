# Capuchin Backend Best Practices

This document outlines the coding standards, patterns, and best practices strictly enforced across the Go backend codebase.

## 1. Centralized Configuration
Environment variables should never be accessed arbitrarily via `os.Getenv` throughout the business logic. 
- All environment variables are loaded, parsed, and validated cleanly within `internal/config`.
- Missing required configurations immediately trigger a `log.Fatal()`, preventing the application from booting into a broken state.

## 2. Interface-Driven Services
Services are defined using Go interfaces.
```go
type TodoService interface {
	GetTodos(userID uuid.UUID) ([]models.Todo, error)
    // ...
}
```
This enables decoupled abstractions. If we decide to swap the database layer out for an ORM or a NoSQL database, we only rewrite the struct that satisfies the interface, leaving the handlers untouched. It also allows for generating mock services for unit testing the handler layer.

## 3. Strong Typing and Struct Binding
We utilize Gin's `ShouldBindJSON` alongside struct tags to strictly map and validate incoming requests before processing them. We refuse requests with an HTTP 400 Bad Request if they violate validation tags (e.g., `binding:"required,min=8"` for passwords).

```go
var reqBody struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

## 4. Centralized Domain Errors
The Service layer does not return HTTP status codes or Gin contexts. Instead, it returns standard Go `error` types defined natively within the package.
```go
var (
	ErrUserExists        = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
```
The Handler layer is responsible for translating these domain errors into the correct semantic HTTP response codes (e.g., 409 Conflict, 401 Unauthorized, 404 Not Found).

## 5. Security Practices
- **Password Hashing:** Passwords are never stored or logged in plain text. We utilize the industry-standard `golang.org/x/crypto/bcrypt` to hash and salt passwords with an appropriate computational cost.
- **JWT Hardening:** The JWT middleware strictly forces the `jwt.WithValidMethods([]string{"HS256"})` and `jwt.WithExpirationRequired()` validators to prevent token tampering or downgrade attacks.
- **Data Isolation:** All protected routes fetch the user ID strictly from the verified JWT token (`c.MustGet("userID")`) injected by the middleware. We never trust `user_id` passed in the HTTP body, effectively preventing lateral data access (IDOR).
- **Graceful Error Recovery:** A global error recovery middleware traps unhandled panics, logs them securely on the server-side, and returns a generic `500 Internal Server Error` to the client, preventing stack trace exposure.

## 6. Resource Management
- **Database Iterator Safety:** When iterating through `rows.Next()`, we explicitly check `rows.Err()` afterward. This catches scenarios where the iteration abruptly halted due to mid-network disconnects or corruption.
- **Background Cleanup:** Dead data (expired logout tokens) is swept away gracefully by an isolated Go routine initialized at startup `go func() { ... }()`, preventing table bloat over time.
