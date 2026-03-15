# Capuchin Backend API Contract

This document provides a comprehensive overview of all available backend endpoints, their expected JSON payloads, requirements, and responses.

## Base URL
When running locally: `http://localhost:8080`

All endpoints return JSON responses. Errors are formatted as `{"error": "description"}`.

---

## Public Endpoints

### 1. Health Check
Checks if the server is running.
- **URL**: `/health`
- **Method**: `GET`
- **Auth Required**: No
- **Response**: `200 OK`
  ```json
  {"status": "ok"}
  ```

### 2. User Signup
Registers a new user account.
- **URL**: `/signup`
- **Method**: `POST`
- **Auth Required**: No
- **Payload**:
  ```json
  {
      "email": "user@example.com", // Required, must be valid email
      "password": "strongpassword123" // Required, min 8 characters
  }
  ```
- **Responses**:
  - `201 Created`: `{"message": "User created successfully"}`
  - `400 Bad Request`: Validation failure (missing fields or password < 8 chars)
  - `409 Conflict`: User with the specified email already exists

### 3. User Login
Authenticates a user and returns a JWT token.
- **URL**: `/login`
- **Method**: `POST`
- **Auth Required**: No
- **Payload**:
  ```json
  {
      "email": "user@example.com",
      "password": "strongpassword123"
  }
  ```
- **Responses**:
  - `200 OK`: `{"token": "ey..."}` (Use this token as a Bearer token in subsequent Protected requests)
  - `401 Unauthorized`: Invalid credentials

---

## Protected Endpoints
All protected endpoints are grouped under `/api/user/`. They require a valid JWT token in the `Authorization` header.
**Header Format:** `Authorization: Bearer <your_jwt_token>`

### 4. Logout User (Revoke Token)
Logs the current user out by adding their JWT to a blacklist.
- **URL**: `/api/user/logout`
- **Method**: `POST`
- **Auth Required**: Yes
- **Responses**:
  - `200 OK`: `{"message": "Logged out successfully"}`
  - `401 Unauthorized`: Token is missing, expired, or already revoked

### 5. Get All Todos
Retrieves all todo items belonging strictly to the authenticated user.
- **URL**: `/api/user/todo`
- **Method**: `GET`
- **Auth Required**: Yes
- **Responses**:
  - `200 OK`: 
    ```json
    [
        {
            "id": "uuid-string",
            "item": "Buy groceries",
            "completed": false
        }
    ]
    ```

### 6. Create Todo
Adds a new todo item for the authenticated user.
- **URL**: `/api/user/todo`
- **Method**: `POST`
- **Auth Required**: Yes
- **Payload**:
  ```json
  {
      "item": "Review pull requests", // Required
      "completed": false // Optional, defaults to false
  }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
        "id": "new-uuid-string",
        "item": "Review pull requests",
        "completed": false
    }
    ```
  - `400 Bad Request`: Missing the required `item` field

### 7. Partially Update Todo (`PATCH`)
Updates specific fields (the text content, the completion status, or both) of an existing todo. It strictly enforces that the todo `id` provided in the path belongs to the authenticated user.
- **URL**: `/api/user/todo/:id` (Replace `:id` with the UUID of the todo)
- **Method**: `PATCH`
- **Auth Required**: Yes
- **Payload**: Provide one or both fields.
  ```json
  {
      "item": "Review 5 pull requests", // Optional
      "completed": true // Optional
  }
  ```
- **Responses**:
  - `200 OK`: Returns the updated todo schema as seen in GET.
  - `400 Bad Request`: Invalid UUID format in URL or invalid JSON
  - `404 Not Found`: Todo does not exist or does not belong to the user

### 8. Delete Todo
Deletes a specific todo belonging strictly to the authenticated user.
- **URL**: `/api/user/todo/:id` (Replace `:id` with the UUID of the todo)
- **Method**: `DELETE`
- **Auth Required**: Yes
- **Responses**:
  - `200 OK`: `{"message": "Todo deleted successfully"}`
  - `400 Bad Request`: Invalid UUID format in URL
  - `404 Not Found`: Todo does not exist or does not belong to the user
