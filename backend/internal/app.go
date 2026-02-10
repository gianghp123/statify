package internal

import (
	"context"
	"log"

	"github.com/gianghp/statify/internal/core/sse"
	"github.com/gianghp/statify/internal/database"
	middlewares "github.com/gianghp/statify/internal/middlewares"
	analyticsModule "github.com/gianghp/statify/internal/modules/analytics"
	"github.com/gianghp/statify/internal/modules/analytics/metrics"
	metricRepo "github.com/gianghp/statify/internal/modules/analytics/repository"
	analyticsService "github.com/gianghp/statify/internal/modules/analytics/service"
	"github.com/gianghp/statify/internal/modules/analytics/wrapper"
	authModule "github.com/gianghp/statify/internal/modules/auth"
	"github.com/gianghp/statify/internal/modules/auth/policy"
	authService "github.com/gianghp/statify/internal/modules/auth/service"
	deploymentModule "github.com/gianghp/statify/internal/modules/deployment"
	deploymentRepository "github.com/gianghp/statify/internal/modules/deployment/repository"
	deploymentService "github.com/gianghp/statify/internal/modules/deployment/service"
	jobQueueRepository "github.com/gianghp/statify/internal/modules/job-queue/repository"
	projectModule "github.com/gianghp/statify/internal/modules/project"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	projectService "github.com/gianghp/statify/internal/modules/project/service"
	uploadSessionRepository "github.com/gianghp/statify/internal/modules/upload-session/repository"
	userModule "github.com/gianghp/statify/internal/modules/user"
	userRepository "github.com/gianghp/statify/internal/modules/user/repository"
	userService "github.com/gianghp/statify/internal/modules/user/service"
	storageMinio "github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *gorm.DB, minioClient *minio.Client, broker *sse.Broker, listener *database.PostgresNotificationListener) *App {
	router := gin.New()

	// Add the colorful logger manually
	router.Use(gin.Logger())

	// Highly recommended: add recovery so your server doesn't die on errors
	router.Use(gin.Recovery())

	//Middlewares
	router.Use(middlewares.ErrorLoggingMiddleware())

	api := router.Group("/api/v1")

	// Repositories
	userRepo := userRepository.NewUserRepository(db)
	projectRepo := projectRepository.NewProjectRepository(db)
	deploymentRepo := deploymentRepository.NewDeploymentRepository(db)
	jobQueueRepo := jobQueueRepository.NewJobQueueRepository(db)
	uploadSessionRepo := uploadSessionRepository.NewUploadSessionRepository(db)

	// Services
	bcryptUtils := bcrypt.NewBcryptUtils()
	jwtService := authService.NewJwtService()

	//Create admin if not exist
	if err := authService.CreateAdmin(userRepo, bcryptUtils); err != nil {
		log.Fatal(err)
	}

	//Metric collector
	metricsRepo := metricRepo.NewDeploymentMetricsMinuteRepository(db)
	metricCollector := metrics.NewMetricCollector(metricsRepo)
	metricCollector.Start(context.Background())
	analyzerWrapper := wrapper.NewAnalyzerWrapper(metricCollector)

	accessPolicy := policy.NewAccessPolicy(projectRepo, deploymentRepo)

	authSvc := authService.NewAuthService(userRepo, bcryptUtils, jwtService)
	minioClientWrapper := storageMinio.NewClient(minioClient)
	projectSvc := projectService.NewProjectService(projectRepo, deploymentRepo, minioClientWrapper, accessPolicy)
	deploymentSvc := deploymentService.NewDeploymentService(deploymentRepo, projectRepo, jobQueueRepo, uploadSessionRepo, minioClientWrapper, accessPolicy)
	userSvc := userService.NewUserService(userRepo)
	analyticsSvc := analyticsService.NewAnalyticsService(metricsRepo)

	// Controllers
	authCtrl := authModule.NewAuthController(authSvc)
	projectCtrl := projectModule.NewProjectController(projectSvc)
	deploymentCtrl := deploymentModule.NewDeploymentController(deploymentSvc, broker)
	userCtrl := userModule.NewUserController(userSvc)
	analyticsCtrl := analyticsModule.NewAnalyticsController(broker, analyticsSvc)

	// Modules
	authMod := authModule.NewAuthModule(authCtrl)
	projectMod := projectModule.NewProjectModule(projectCtrl)
	deploymentMod := deploymentModule.NewDeploymentModule(deploymentCtrl)
	userMod := userModule.NewUserModule(userCtrl)
	analyticsMod := analyticsModule.NewAnalyticsModule(analyticsCtrl)

	//Middlewares
	authMiddleware := middlewares.AuthMiddleware(jwtService)
	// roleMiddleware := middlewares.RoleMiddleware(enums.UserRoleAdmin)

	// Register Routes
	authMod.RegisterRoutes(api, authMiddleware)
	projectMod.RegisterRoutes(api, authMiddleware)
	deploymentMod.RegisterRoutes(api, authMiddleware, analyzerWrapper)
	userMod.RegisterRoutes(api, authMiddleware)
	analyticsMod.RegisterRoutes(api, authMiddleware)

	// Start the listener in a goroutine
	go listener.Run(context.Background(), sse.DeploymentStatusEvent)
	go listener.Run(context.Background(), sse.AnalyticsEvent)

	return &App{
		Router: router,
	}
}
