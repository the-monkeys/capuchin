package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// create a default gin router
	r := gin.Default()

	// create a Get route at "/ping"
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server started",
		})
	})

	//start the server on localhost:8080
	r.Run()
}
