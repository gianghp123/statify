package user

import (
	"github.com/gianghp/statify/internal/modules/user/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// TODO: Implement
}
