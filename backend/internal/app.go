package internal

import (
	authModule "github.com/gianghp/statify/internal/modules/auth"
	authService "github.com/gianghp/statify/internal/modules/auth/service"
	deploymentModule "github.com/gianghp/statify/internal/modules/deployment"
	deploymentRepository "github.com/gianghp/statify/internal/modules/deployment/repository"
	projectModule "github.com/gianghp/statify/internal/modules/project"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	userRepoitory "github.com/gianghp/statify/internal/modules/user/repository"
	"github.com/gianghp/statify/internal/utils/bcrypt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *gorm.DB) *App {
	router := gin.New()
	api := router.Group("/api/v1")

	bcryptUtils := bcrypt.NewBcryptUtils()
	jwtService := authService.NewJwtService()

	userRepo := userRepoitory.NewUserRepository(db)
	projectRepo := projectRepository.NewProjectRepository(db)
	deploymentRepo := deploymentRepository.NewDeploymentRepository(db)

	authMod := authModule.NewAuthModule(userRepo, bcryptUtils, jwtService)
	authMod.RegisterRoutes(api)

	projectMod := projectModule.NewProjectModule(projectRepo)
	projectMod.RegisterRoutes(api)

	deploymentMod := deploymentModule.NewDeploymentModule(deploymentRepo)
	deploymentMod.RegisterRoutes(api)

	return &App{
		Router: router,
	}
}
