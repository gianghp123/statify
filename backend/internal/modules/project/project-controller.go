package project

import (
	"net/http"
	"strconv"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
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

	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	projects, err := c.service.ListProjects(ctx, userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewPaginatedApiResponse(http.StatusOK, "Projects retrieved successfully", projects.Entities, &projects.Pagination))
}

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	var req request.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.HandleApiError(ctx, core.BadRequestError(err.Error()))
		return
	}

	projectDto, err := c.service.CreateProject(ctx, userID, &req)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, core.NewApiResponse(http.StatusCreated, "Project created successfully", projectDto))
}

func (c *ProjectController) GetProject(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)

	if err != nil {
		core.HandleApiError(ctx, core.UnauthorizedError())
		return
	}

	idStr := ctx.Param("project_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		core.HandleApiError(ctx, core.BadRequestError("Invalid project ID"))
		return
	}

	projectDto, err := c.service.GetProjectByID(ctx, uint(id), userID)
	if err != nil {
		core.HandleApiError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, core.NewApiResponse(http.StatusOK, "Project retrieved successfully", projectDto))
}
