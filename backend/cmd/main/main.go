package main

import (
	"log"
	"net/http"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/gianghp/statify/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const swaggerUIHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.52.0/swagger-ui.css" >
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.52.0/swagger-ui-bundle.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.52.0/swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
      // Initialize Swagger UI
      const ui = SwaggerUIBundle({
        // **This URL must match the Gin route that serves your swagger.yaml**
        url: "/api/swagger/spec", 
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      })
    }
    </script>
</body>
</html>
`

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, using system environment variables")
	}

	_, err := database.InitDatabase()

	if err != nil {
		log.Fatal("Failed when connecting to the database")
		return
	}

	log.Println("Database connected successfully")
	r := gin.Default()

	// Setup Swagger documentation

	// Setup all routes
	// modules.SetupAllRoutes(r, db)

	r.GET("/api/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})

	// 2. Serve the OpenAPI 3.0 specification file statically
	r.StaticFile("/api/swagger/spec", "./docs/swagger.yaml")
	// Start the server
	r.Run(":8000") // Listen and serve on 0.0.0.0:8000

}
