package auth

import (
	"net/http"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/auth/dtos/request"
	"github.com/gianghp/statify/internal/modules/auth/service"
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

	userDto, err := c.service.Register(ctx, &request)

	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "User registered successfully", userDto))
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	result, err := c.service.Login(ctx, &req)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Login successful", result))
}

func (c *AuthController) GitHubLogin(ctx *gin.Context) {
	// TODO: Implement GitHub OAuth initiation
}

func (c *AuthController) GitHubCallback(ctx *gin.Context) {
	// TODO: Implement GitHub OAuth callback
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	userDto, err := c.service.Me(ctx, userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "User profile retrieved successfully", userDto))
}
