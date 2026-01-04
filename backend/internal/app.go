package internal

import (
	"log"

	middlewares "github.com/gianghp/statify/internal/middleware"
	authModule "github.com/gianghp/statify/internal/modules/auth"
	authService "github.com/gianghp/statify/internal/modules/auth/service"
	deploymentModule "github.com/gianghp/statify/internal/modules/deployment"
	deploymentRepository "github.com/gianghp/statify/internal/modules/deployment/repository"
	deploymentService "github.com/gianghp/statify/internal/modules/deployment/service"
	projectModule "github.com/gianghp/statify/internal/modules/project"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	projectService "github.com/gianghp/statify/internal/modules/project/service"
	userModule "github.com/gianghp/statify/internal/modules/user"
	userRepository "github.com/gianghp/statify/internal/modules/user/repository"
	userService "github.com/gianghp/statify/internal/modules/user/service"
	"github.com/gianghp/statify/internal/utils/bcrypt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *gorm.DB) *App {
	router := gin.New()

	// Add the colorful logger manually
	router.Use(gin.Logger())

	// Highly recommended: add recovery so your server doesn't die on errors
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	// Repositories
	userRepo := userRepository.NewUserRepository(db)
	projectRepo := projectRepository.NewProjectRepository(db)
	deploymentRepo := deploymentRepository.NewDeploymentRepository(db)

	// Services
	bcryptUtils := bcrypt.NewBcryptUtils()
	jwtService := authService.NewJwtService()

	//Create admin if not exist
	if err := authService.CreateAdmin(userRepo, bcryptUtils); err != nil {
		log.Fatal(err)
	}

	authSvc := authService.NewAuthService(userRepo, bcryptUtils, jwtService)
	projectSvc := projectService.NewProjectService(projectRepo, deploymentRepo)
	deploymentSvc := deploymentService.NewDeploymentService(deploymentRepo, projectRepo)
	userSvc := userService.NewUserService(userRepo)

	// Controllers
	authCtrl := authModule.NewAuthController(authSvc)
	projectCtrl := projectModule.NewProjectController(projectSvc)
	deploymentCtrl := deploymentModule.NewDeploymentController(deploymentSvc)
	userCtrl := userModule.NewUserController(userSvc)

	// Modules
	authMod := authModule.NewAuthModule(authCtrl)
	projectMod := projectModule.NewProjectModule(projectCtrl)
	deploymentMod := deploymentModule.NewDeploymentModule(deploymentCtrl)
	userMod := userModule.NewUserModule(userCtrl)

	//Middlewares
	authMiddleware := middlewares.AuthMiddleware(jwtService)
	// roleMiddleware := middlewares.RoleMiddleware(enums.UserRoleAdmin)

	// Register Routes
	authMod.RegisterRoutes(api, authMiddleware)
	projectMod.RegisterRoutes(api, authMiddleware)
	deploymentMod.RegisterRoutes(api, authMiddleware)
	userMod.RegisterRoutes(api, authMiddleware)

	return &App{
		Router: router,
	}
}
