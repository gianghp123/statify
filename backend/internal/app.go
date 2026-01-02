package internal

import (
	authModule "github.com/gianghp/statify/internal/modules/auth"
	authService "github.com/gianghp/statify/internal/modules/auth/service"
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

	userRepo := userRepoitory.NewUserRepository(db)
	bcryptUtils := bcrypt.NewBcryptUtils()
	jwtService := authService.NewJwtService()

	authModule := authModule.NewAuthModule(userRepo, bcryptUtils, jwtService)
	authModule.RegisterRoutes(api)

	return &App{
		Router: router,
	}
}
