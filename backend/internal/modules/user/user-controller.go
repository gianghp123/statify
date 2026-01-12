package user

import (
	"net/http"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/user/dtos/request"
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
	// userID, err := utils.GetUserIDFromContext(ctx)
	// if err != nil {
	// 	core.HandleApiError(ctx, core.UnauthorizedError())
	// 	return
	// }

	var req request.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	// TODO: Implement actual update logic with service
	ctx.JSON(http.StatusOK, core.NewApiResponse[any](http.StatusOK, "Profile updated successfully (placeholder)", nil))
}

func (c *UserController) GetUser(ctx *gin.Context) {
	// TODO: Implement actual get logic with service
	ctx.JSON(http.StatusOK, core.NewApiResponse[any](http.StatusOK, "User details retrieved successfully (placeholder)", nil))
}

func (c *UserController) DeleteAccount(ctx *gin.Context) {
	// TODO: Implement actual delete logic with service
	ctx.JSON(http.StatusOK, core.NewApiResponse[any](http.StatusOK, "Account deleted successfully (placeholder)", nil))
}
