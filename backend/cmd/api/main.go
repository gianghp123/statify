package main

import (
	"log"
	"net/http"

	"github.com/gianghp/statify/internal"
	"github.com/gianghp/statify/internal/assets"
	"github.com/gianghp/statify/internal/configs"
	"github.com/gianghp/statify/internal/core/sse"
	"github.com/gianghp/statify/internal/database"
	"github.com/gianghp/statify/internal/database/migrations"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, using system environment variables")
	}

	connectStr := configs.LoadDatabaseConfig()

	db, err := database.InitDatabase(connectStr)

	if err != nil {
		log.Fatal("Failed when connecting to the database")
		return
	}

	log.Println("Database connected successfully")

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDB, "."); err != nil {
		panic(err)
	}

	accessKeyID, secretAccessKey := configs.LoadMinioConfig()
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(utils.GetEnv("MINIO_URL", "localhost:9000"), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal("Failed to initialize minio client", err)
		return
	}

	broker := sse.NewBroker()

	listener := database.NewPostgresNotificationListener(connectStr, broker)

	app := internal.NewApp(db, minioClient, broker, listener)

	app.Router.Use(cors.Default())

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
