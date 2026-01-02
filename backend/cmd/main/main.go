package main

import (
	"log"
	"net/http"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/gianghp/statify/internal"
	"github.com/gianghp/statify/internal/assets"
	"github.com/gianghp/statify/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, using system environment variables")
	}

	db, err := database.InitDatabase()

	if err != nil {
		log.Fatal("Failed when connecting to the database")
		return
	}

	log.Println("Database connected successfully")

	app := internal.NewApp(db)

	app.Router.GET("/api/v1/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(assets.SwaggerUIHTML))
	})

	// 2. Serve the OpenAPI 3.0 specification file statically
	app.Router.GET("/api/v1/swagger/spec", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/yaml; charset=utf-8", assets.OpenAPISpec)
	})
	// Start the server
	app.Router.Run(":8000")

}
