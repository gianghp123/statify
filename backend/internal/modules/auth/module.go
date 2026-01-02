package auth

import (
	"github.com/gianghp/statify/internal/modules/auth/service"
	"github.com/gianghp/statify/internal/modules/user/repository"
	"github.com/gianghp/statify/internal/utils/bcrypt"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	controller *AuthController
}

func NewAuthModule(userRepo repository.IUserRepository, bcryptUtils bcrypt.IBcryptUtils, jwtService service.IJwtService) *AuthModule {
	authService := service.NewAuthService(userRepo, bcryptUtils, jwtService)

	return &AuthModule{
		controller: NewAuthController(authService),
	}
}

func (m *AuthModule) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", m.controller.Register)
		auth.POST("/login", m.controller.Login)
		auth.GET("/github", m.controller.GitHubLogin)
		auth.GET("/github/callback", m.controller.GitHubCallback)
		auth.GET("/me", m.controller.Me)
	}
}
