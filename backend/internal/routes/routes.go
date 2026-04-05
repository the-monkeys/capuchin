package routes

import (
	"capuchin/internal/handlers"
	"capuchin/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler) {

	router.Use(middleware.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	router.GET("/api-docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	router.POST("/signup", authHandler.Signup)
	router.POST("/login", authHandler.Login)

	// Group authenticated routes so auth middleware is applied consistently.
	protected := router.Group("/api/user")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/todo", todoHandler.GetTodos)
		protected.POST("/todo", todoHandler.AddTodo)
		protected.PATCH("/todo/:id", todoHandler.UpdateTodo)
		protected.DELETE("/todo/:id", todoHandler.DeleteTodo)
	}
}

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Capuchin API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/openapi.yaml",
      dom_id: "#swagger-ui",
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis],
      layout: "BaseLayout"
    });
  </script>
</body>
</html>`
