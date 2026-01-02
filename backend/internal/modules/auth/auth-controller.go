package auth

import (
	"net/http"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/auth/dtos/request"
	"github.com/gianghp/statify/internal/modules/auth/service"
	"github.com/gianghp/statify/internal/modules/user/dtos/response"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var request request.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	user, err := c.service.Register(&request)

	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	userDto, err := utils.EntityToDto[response.UserDto](user)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "User registered successfully", userDto))
}

func (c *AuthController) Login(ctx *gin.Context) {
	// TODO: Implement
}

func (c *AuthController) GitHubLogin(ctx *gin.Context) {
	// TODO: Implement GitHub OAuth initiation
}

func (c *AuthController) GitHubCallback(ctx *gin.Context) {
	// TODO: Implement GitHub OAuth callback
}

func (c *AuthController) Me(ctx *gin.Context) {
	// TODO: Implement
}
