package handlers

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"capuchin/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
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
	u.ID = uuid.New()
	u.Email = req.Email
	u.PasswordHash = string(hash)

	_, err := database.DB.Exec("INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)", u.ID, u.Email, u.PasswordHash)
	if err != nil {
		c.JSON(500, gin.H{"error": "User already exists or db error"})
		return
	}
	c.JSON(201, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var u models.User
	err := database.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", req.Email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID.String(),
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(config.JWTKey)
	c.JSON(200, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	// The auth middleware already verified the token, so we just need to get the raw token string
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(400, gin.H{"error": "Authorization header missing"})
		return
	}
	
	// Remove "Bearer " prefix
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	// Calculate expiration based on JWT claims
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})

	var expTime time.Time
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expTime = time.Unix(int64(exp), 0)
		} else {
			expTime = time.Now().Add(72 * time.Hour) // Fallback
		}
	} else {
		c.JSON(400, gin.H{"error": "Invalid token components"})
		return
	}

	// Insert into blacklisted_tokens table
	_, err := database.DB.Exec("INSERT INTO blacklisted_tokens (token, expired_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING", tokenStr, expTime)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}
