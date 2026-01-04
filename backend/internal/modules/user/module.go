package user

import (
	"github.com/gin-gonic/gin"
)

type UserModule struct {
	controller *UserController
}

func NewUserModule(controller *UserController) *UserModule {
	return &UserModule{
		controller: controller,
	}
}

func (m *UserModule) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	user := rg.Group("/users")
	{
		user.PUT("/profile", authMiddleware, m.controller.UpdateProfile)
	}
}
