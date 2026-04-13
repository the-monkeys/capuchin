package middleware

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		// Accept standard Authorization header format without forcing clients to preprocess it.
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		var exists bool
		// Check revocation before claim extraction so logout takes effect immediately.
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM blacklisted_tokens WHERE token=$1)", tokenStr).Scan(&exists)
		if err == nil && exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token has been revoked"})
			return
		}

		// Restrict acceptable algorithms and require exp to reduce token confusion attacks.
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// JWT numbers are decoded as float64 by encoding/json.
			raw, ok := claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
				return
			}

			c.Set("userID", int64(raw))
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}
