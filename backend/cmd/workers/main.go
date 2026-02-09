package main

import (
	"context"
	"log"
	"time"

	"github.com/gianghp/statify/internal/configs"
	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/database"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/deployment/service"
	"github.com/gianghp/statify/internal/modules/deployment/workers"
	jobQueueRepo "github.com/gianghp/statify/internal/modules/job-queue/repository"
	storageMinio "github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	deploymentRepository := repository.NewDeploymentRepository(db)
	minioClientWrapper := storageMinio.NewClient(minioClient)
	fileProcessor := service.NewFileProcessor(minioClientWrapper, deploymentRepository, utils.GetEnv("MINIO_BUCKET", "statify"))
	jobQueueRepository := jobQueueRepo.NewJobQueueRepository(db)
	deploymentWorker := workers.NewDeploymentWorker(jobQueueRepository, deploymentRepository, fileProcessor)

	deploymentWorker.Run(context.Background(), core.MaxProcessGoroutines, core.MaxDeleteGoroutines, 2*time.Second, 5)

}
