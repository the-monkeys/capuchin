package middleware

import (
	"capuchin/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DBHealthCheck returns a Gin middleware that responds 503 Service Unavailable
// when the database is not reachable, preventing handlers from executing
// against a nil or unhealthy DB connection.
func DBHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !database.IsHealthy() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "database unavailable",
			})
			return
		}
		c.Next()
	}
}
