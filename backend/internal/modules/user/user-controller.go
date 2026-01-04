package user

import (
	"github.com/gianghp/statify/internal/modules/user/service"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	_, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		// core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}
	// TODO: Implement actual update logic with service
}
