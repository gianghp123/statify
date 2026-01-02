package user

import "github.com/gin-gonic/gin"

type UserRoutes struct {
	controller *UserController
}

func NewUserRoutes(controller *UserController) *UserRoutes {
	return &UserRoutes{controller: controller}
}

func (r *UserRoutes) SetupUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/users")
	{
		user.PUT("/profile", r.controller.UpdateProfile)
	}
}
