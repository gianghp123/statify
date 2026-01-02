package project

import (
	"net/http"
	"strconv"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/dtos/response"
	"github.com/gianghp/statify/internal/modules/project/service"
	"github.com/gianghp/statify/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	service *service.ProjectService
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

func (c *ProjectController) ListProjects(ctx *gin.Context) {
	// For MVP, we assume user ID is 1 or get from context if auth middleware is there
	// TODO: Get user ID from context
	userID := uint(1)

	projects, err := c.service.ListProjects(userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	projectDtos, err := utils.EntitiesToDto[response.ProjectDto](projects.Entities)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Projects retrieved successfully", projectDtos, &projects.Pagination))
}

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req request.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	// TODO: Get user ID from context
	userID := uint(1)

	project, err := c.service.CreateProject(userID, &req)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	projectDto, err := utils.EntityToDto[response.ProjectDto](project)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "Project created successfully", projectDto))
}

func (c *ProjectController) GetProject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	project, err := c.service.GetProjectByID(uint(id))
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	if project == nil {
		core.HandleApiError(ctx, core.NotFoundError())
		return
	}

	projectDto, err := utils.EntityToDto[response.ProjectDto](project)
	if err != nil {
		core.HandleApiError(ctx, core.InternalError())
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Project retrieved successfully", projectDto))
}
