package deployment

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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

func (c *DeploymentController) GetGlobalDeploymentHistory(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	page, limit := utils.GetPaginationConfig(ctx)
	deployments, err := c.service.GetGlobalDeploymentHistory(ctx, userID, page, limit)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Global deployment history retrieved successfully", deployments.Entities, &deployments.Pagination))
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

	page, limit := utils.GetPaginationConfig(ctx)
	deployments, err := c.service.GetHistory(ctx, userID, uint(projectID), page, limit)
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

func (c *DeploymentController) ServeFiles(ctx *gin.Context) {
	host := strings.Split(ctx.Request.Host, ".")[0]

	clientEtag := ctx.GetHeader("If-None-Match")
	fileName := ctx.Param("file_name")

	if fileName == "" || fileName == "/" {
		fileName = "index.html"
	}

	// Optional: strip the leading slash if it exists
	fileName = strings.TrimPrefix(fileName, "/")

	log.Println("host", host, "fileName", fileName, "clientEtag", clientEtag)
	fileDTO, err := c.service.GetCurrentDeploymentFilesByProjectSubdomain(ctx, host, fileName, clientEtag)

	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	if fileDTO.Stream != nil {
		defer fileDTO.Stream.Close()
	}

	if fileDTO.NotModified {
		ctx.Status(http.StatusNotModified)
		return
	}

	ctx.DataFromReader(http.StatusOK, fileDTO.Size, fileDTO.ContentType, fileDTO.Stream, fileDTO.Headers)
}

func (c *DeploymentController) DeleteDeployment(ctx *gin.Context) {
	// userID, err := utils.GetUserIDFromContext(ctx)
	// if err != nil {
	// 	core.HandleApiError(ctx, core.UnauthorizedError())
	// 	return
	// }

	// TODO: Implement actual delete logic with service
	ctx.JSON(http.StatusOK, core.NewApiResponse[any](http.StatusOK, "Deployment deleted successfully (placeholder)", nil))
}

func (c *DeploymentController) TurnDeploymentLive(ctx *gin.Context) {
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

	err = c.service.TurnDeploymentLive(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, core.NewApiResponse[any](http.StatusNoContent, "Deployment turned live successfully", nil))
}

func (c *DeploymentController) TurnDeploymentOffline(ctx *gin.Context) {
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

	err = c.service.TurnDeploymentOffline(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, core.NewApiResponse[any](http.StatusNoContent, "Deployment turned offline successfully", nil))
}
