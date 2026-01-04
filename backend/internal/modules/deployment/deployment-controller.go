package deployment

import (
	"net/http"
	"strconv"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/service"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-gonic/gin"
)

type DeploymentController struct {
	service *service.DeploymentService
}

func NewDeploymentController(service *service.DeploymentService) *DeploymentController {
	return &DeploymentController{service: service}
}

func (c *DeploymentController) GetHistory(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	deployments, err := c.service.GetHistory(ctx, userID, uint(projectID))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Deployment history retrieved successfully", deployments.Entities, &deployments.Pagination))
}

func (c *DeploymentController) CreateDeployment(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	var req request.CreateDeploymentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	deploymentDto, err := c.service.CreateDeployment(ctx, userID, uint(projectID), &req)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "Deployment created successfully", deploymentDto))
}

func (c *DeploymentController) GetStatus(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	deploymentDto, err := c.service.GetDeploymentByID(ctx, userID, uint(id))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Deployment status retrieved successfully", deploymentDto))
}
