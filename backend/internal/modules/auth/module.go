package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	controller *AuthController
}

func NewAuthModule(controller *AuthController) *AuthModule {
	return &AuthModule{
		controller: controller,
	}
}

func (m *AuthModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", m.controller.Register)
		auth.POST("/login", m.controller.Login)
		auth.GET("/github", m.controller.GitHubLogin)
		auth.GET("/github/callback", m.controller.GitHubCallback)
		auth.GET("/me", authMiddleware, m.controller.Me)
	}
}
