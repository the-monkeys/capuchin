package middleware

import (
	"capuchin/internal/config"
	"capuchin/internal/database"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthRequired validates JWT token and injects user ID into context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}
		
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Check if token is blacklisted
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM blacklisted_tokens WHERE token=$1)", tokenStr).Scan(&exists)
		if err == nil && exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token has been revoked"})
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check expiration manually if needed, though jwt.Parse handles "exp" standard claim
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.AbortWithStatusJSON(401, gin.H{"error": "Token has expired"})
					return
				}
			}

			// Extract user ID
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
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}
