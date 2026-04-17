package middleware

import (
	"capuchin/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DBHealthCheck returns 503 Service Unavailable when the database is unreachable.
// This prevents requests from reaching handlers that depend on a live DB connection.
func DBHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !database.IsHealthy() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "service temporarily unavailable",
			})
			return
		}
		c.Next()
	}
}
