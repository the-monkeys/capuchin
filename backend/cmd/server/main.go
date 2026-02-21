package main

import (
	"capuchin/internal/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func main() {
	if len(jwtKey) == 0 {
		jwtKey = []byte("secret") // Default for dev if env not set
	}

	// Connect to Postgres
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	initDB()

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)

	// Protected Routes
	user := r.Group("/user")
	user.Use(authMiddleware())
	{
		user.GET("/todos", getTodos)
		user.POST("/todos", addTodo)
		user.PATCH("/todos/:id", toggleTodo)
		user.PATCH("/todos/:id/edit", editTodo)
		user.DELETE("/todos/:id", deleteTodo)
	}

	r.Run(":8080")
}

func initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		item TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		user_id INTEGER REFERENCES users(id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
}

// --- Auth Handlers ---

func signupHandler(c *gin.Context) {
	var u models.User
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u.Email = req.Email
	u.PasswordHash = string(hash)

	err := db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", u.Email, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "User already exists or db error"})
		return
	}
	c.JSON(201, gin.H{"message": "User created"})
}

func loginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var u models.User
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", req.Email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(200, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatus(401)
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", int64(claims["user_id"].(float64)))
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}

// --- Todo Handlers ---

func getTodos(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	rows, err := db.Query("SELECT id, item, completed FROM todos WHERE user_id=$1", userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
			continue
		}
		todos = append(todos, t)
	}
	c.JSON(200, todos)
}

func addTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var t models.Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	t.ID = uuid.New().String()
	t.UserID = userID

	_, err := db.Exec("INSERT INTO todos (id, item, completed, user_id) VALUES ($1, $2, $3, $4)", t.ID, t.Item, t.Completed, t.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, t)
}

func toggleTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id := c.Param("id")

	var t models.Todo
	// Toggle and return new state
	err := db.QueryRow(`
		UPDATE todos SET completed = NOT completed 
		WHERE id=$1 AND user_id=$2 
		RETURNING id, item, completed`, id, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(200, t)
}

func editTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id := c.Param("id")
	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var t models.Todo
	err := db.QueryRow(`
		UPDATE todos SET item=$1 
		WHERE id=$2 AND user_id=$3 
		RETURNING id, item, completed`, req.Item, id, userID).Scan(&t.ID, &t.Item, &t.Completed)

	if err != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(200, t)
}

func deleteTodo(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	id := c.Param("id")
	db.Exec("DELETE FROM todos WHERE id=$1 AND user_id=$2", id, userID)
	c.JSON(200, gin.H{"message": "deleted"})
}
