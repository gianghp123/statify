package deployment

import (
	"net/http"
	"strconv"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/response"
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
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	deployments, err := c.service.GetHistory(uint(projectID))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	deploymentDtos, err := utils.EntitiesToDto[response.DeploymentDto](deployments.Entities)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Deployment history retrieved successfully", deploymentDtos, &deployments.Pagination))
}

func (c *DeploymentController) CreateDeployment(ctx *gin.Context) {
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

	deployment, err := c.service.CreateDeployment(uint(projectID), &req)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	deploymentDto, err := utils.EntityToDto[response.DeploymentDto](deployment)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "Deployment created successfully", deploymentDto))
}

func (c *DeploymentController) GetStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid deployment ID"))
		return
	}

	deployment, err := c.service.GetDeploymentByID(uint(id))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	if deployment == nil {
		core.HandleApiError(ctx, core.NotFoundError())
		return
	}

	deploymentDto, err := utils.EntityToDto[response.DeploymentDto](deployment)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Deployment status retrieved successfully", deploymentDto))
}
