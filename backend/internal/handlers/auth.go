package handlers

import (
	"capuchin/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{authService: svc}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var reqBody struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: missing fields or invalid format. Password must be >= 8 characters."})
		return
	}

	_, err := h.authService.Signup(reqBody.Email, reqBody.Password)
	if err != nil {
		if err == services.ErrUserExists {
			c.JSON(409, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := h.authService.Login(reqBody.Email, reqBody.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(400, gin.H{"error": "Authorization header missing"})
		return
	}

	err := h.authService.Logout(tokenStr)
	if err != nil {
		if err == services.ErrInvalidToken {
			c.JSON(400, gin.H{"error": "Invalid token components"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}
