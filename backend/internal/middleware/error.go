package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler provides a basic recovery and error formatting middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"success": false,
					"details": err,
				})
			}
		}()

		c.Next()

		// If there are errors attached to the context, we can format them
		if len(c.Errors) > 0 {
			// You can implement custom logic based on error types here
			// For now, we'll just log them and ensure a response is sent if not already
			for _, e := range c.Errors {
				log.Printf("Request error: %v", e.Err)
			}

			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   c.Errors.Last().Error(),
					"success": false,
				})
			}
		}
	}
}
