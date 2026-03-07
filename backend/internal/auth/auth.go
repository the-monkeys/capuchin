package auth

import (
	"capuchin/config"
	"capuchin/internal/database"
	"capuchin/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func Init() {
	jwtKey = config.JWTKey
}

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
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(200, gin.H{"token": tokenString})
}

func Middleware() gin.HandlerFunc {
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
			raw, ok := claims["user_id"].(string)
			if !ok {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
				return
			}
			uid, err := uuid.Parse(raw)
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid user id in token"})
				return
			}
			c.Set("userID", uid)
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}
